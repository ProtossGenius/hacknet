package com.yqg.recall.core.utils;


import com.google.common.util.concurrent.ThreadFactoryBuilder;
import com.yqg.recall.common.util.WorkGroup;
import com.yqg.recall.common.util.pipeline.IStore;
import com.yqg.recall.common.util.pipeline.Pipeline;
import com.yqg.tracing.executorservice.ExecutorServiceMdcTraceIdWrapper;
import com.yqg.tracing.executorservice.ExecutorServiceTraceContextWrapper;
import lombok.Data;
import lombok.extern.slf4j.Slf4j;
import org.junit.Test;

import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.LinkedBlockingDeque;
import java.util.concurrent.ThreadPoolExecutor;
import java.util.concurrent.TimeUnit;

@Slf4j
public class PipelineTest {
  private static final ExecutorService EXECUTOR = ExecutorServiceTraceContextWrapper.wrap(
      ExecutorServiceMdcTraceIdWrapper.wrap(
          new ThreadPoolExecutor(
              1,
              200,
              0,
              TimeUnit.MILLISECONDS,
              new LinkedBlockingDeque<>(),
              new ThreadFactoryBuilder().setNameFormat("StrategyRun-pool-%d").build(),
              new ThreadPoolExecutor.CallerRunsPolicy()),
          "trace_id"
      )
  );

  private static void sleep(long ms) {
    try {
      Thread.sleep(ms);
    } catch (Exception e) {
      log.error("sleep error", e);
    }
  }

  public void testNormal() {
    long start = System.nanoTime();
    try {
      /* 许多方法中的每个耗时步骤都是相互独立的，但是用常规写法很难进行优化
      其次，也很难了解到底哪里耗时比较多
      常规写法理论耗时为 68ms
      理论最少耗时为 （ABCD并行+提交）35ms
      实测结果见test()的注释
       */
      Result res = new Result();
      res.setA(getA());//10
      res.setB(getB());//15
      res.setC(getC());//20
      List<String> l = getList();//5
      WorkGroup<Integer> wg = new WorkGroup<>(EXECUTOR);
      l.forEach(it -> wg.add(it, () -> this.getD(it)));// 微> 3
      res.setD(wg.waitAllFinishAndGetResult(log::error));
      submit(res);// 15
    } finally {
      log.info("testNormal 's time cost = " + (System.nanoTime() - start) / 1e6 + " ms");
    }
  }

  public Pipeline getPipeLine() {
    boolean calcTime = false; // 是否打印流程耗时
    // 子流程获得 result.D
    final Pipeline<?, List<Integer>> getD = new Pipeline<>("get d", calcTime, exe -> getList()) // 获得List
        .thenBatch("batch", Iterable::forEach, log::error, (e, list, str, s, cb) -> cb.accept(getD((String) str))); // 批量处理List，转化为需要的D
    // 主流程
    Pipeline pipeline = new Pipeline<>("get list", calcTime, executorService -> -1)
        .thenMerge("merge", Result::new, log::error,
            builder -> builder
                .and("set A",
                    (es, list, sto, result) -> result.setA(getA()))
                .and("get B",
                    (es, list, sto, result) -> result.setB(getB()))
                .and("get C",
                    (es, list, sto, result) -> result.setC(getC()))
                .and("get D",
                    (es, list, sto, result, callbackGroup) -> getD.execute(
                        es,
                        sto,
                        nums -> {
                          result.setD(nums);
                          callbackGroup.unsafeFinishOneCount();
                        }
                    )
                )
        )
        .then("submit", (e, result, s) -> { // 提交
          this.submit(result);
          return 1;
        });

    return pipeline;
  }

  public void testPipeline(Pipeline pipeline) {
    long start = System.nanoTime();
    try {
      pipeline.blockExecute(EXECUTOR, -1, null);
    } finally {
      log.info("testPipeline 's time cost = " + (System.nanoTime() - start) / 1e6 + " ms");
    }
  }

  @Test
  public void test() {
    testNormal();
    testNormal();
    Pipeline pipeline = getPipeLine();
    testPipeline(pipeline);
    testPipeline(pipeline);
    testPipeline(pipeline);
    testPipeline(pipeline);
    testPipeline(pipeline);
    testPipeline(pipeline);
    testPipeline(pipeline);
    /*
19:02:59.743 [main] INFO com.yqg.recall.core.utils.PipelineTest - testNormal 's time cost = 88.404896 ms
19:02:59.822 [main] INFO com.yqg.recall.core.utils.PipelineTest - testNormal 's time cost = 75.38075 ms
19:02:59.880 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 48.621298 ms
19:02:59.916 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 36.296537 ms
19:02:59.956 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 39.804971 ms
19:02:59.994 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 38.142821 ms
19:03:00.037 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 42.204889 ms
19:03:00.077 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 40.475811 ms
19:03:00.113 [main] INFO com.yqg.recall.core.utils.PipelineTest - testPipeline 's time cost = 36.04097 ms
     */
  }

  @Test
  public void testAsFunc() {
    Pipeline pipeline = new Pipeline<>("test get value", (executorService, i, s) -> "input is " + i.toString())
        .then("exception", (executorService, str, s) -> {
          throw new RuntimeException("hh");
        })
        .then("add", (executorService, str, s) -> str + ((ValueStore) s).getEnd())
        .thenMerge("merge", () -> 15, log::error,
            (executorService, s, store, o) -> log.info("a"),
            (executorService, s, store, o) -> log.info("b")
        )
        .then("ex2", (executorService, integer, store) -> {
          throw new RuntimeException("??");
        });
    //ValueStore可以贯穿整个流程
    log.info((String) pipeline.blockExecute(EXECUTOR, 15, new ValueStore()));
    //    19:32:52.568 [main] INFO com.yqg.recall.core.utils.PipelineTest - input is 15!
  }

  private String getA() {
    sleep(10);
    return "hello";
  }

  private String getB() {
    sleep(15);
    return "world";
  }

  private String getC() {
    sleep(20);
    return ".";
  }

  private List<String> getList() {
    sleep(5);
    List<String> l = new ArrayList<>();
    for (int i = 0; i < 15; i++) {
      l.add(Integer.toString(i));
    }

    return l;
  }

  public Integer getD(String o) {
    sleep(3);
    return Integer.parseInt(o);
  }

  void assertEquals(Object a, Object b) {
    if (!a.equals(b)) {
      throw new RuntimeException("<" + a + "> not equals <" + b + ">");
    }
  }

  public void submit(Result result) {
    sleep(15);
    assertEquals(result.getA(), "hello");
    assertEquals(result.getB(), "world");
    assertEquals(result.getC(), ".");
    List<Integer> list = new ArrayList<>();
    for (int i = 0; i < 15; ++i) {
      list.add(i);
    }
    assertEquals(result.getD(), list);
  }

  @Data
  static class Result {
    String a;
    String b;
    String c;
    List<Integer> d;
  }

  @Data
  public static class ValueStore implements IStore {
    String end = "!";
  }

}
