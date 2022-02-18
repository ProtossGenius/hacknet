package com.yqg.recall.common.util.pipeline;

import com.yqg.recall.common.util.WorkGroup;
import lombok.extern.slf4j.Slf4j;

import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 * 流程
 *
 * @param <In>
 * @param <Out>
 */
@Slf4j
public class Pipeline<In, Out> implements IPipeline<In, Out> {
  String name;
  IPipeline<In, Out> action;
  IPipeline previous;
  boolean calcTime = false;

  public Pipeline(String name, boolean calcTime, Function<ExecutorService, Out> action) {
    this.calcTime = calcTime;
    this.name = name;
    this.action = (executorService, in, store) -> action.apply(executorService);
  }

  public Pipeline(String name, Function<ExecutorService, Out> action) {
    this(name, false, action);
  }

  public Pipeline(String name, boolean calcTime, IPipeline<In, Out> action) {
    this.calcTime = calcTime;
    this.name = name;
    this.action = action;
  }

  public Pipeline(String name, IPipeline<In, Out> action) {
    this(name, false, action);
  }

  private Pipeline(String name, IPipeline previous, boolean calcTime, IPipeline<In, Out> action) {
    this.name = name;
    this.previous = previous;
    this.calcTime = calcTime;
    this.action = action;
  }

  // 进入下一个流程
  public <NewOut> Pipeline<Out, NewOut> then(String name, IPipeline<Out, NewOut> action) {
    return new Pipeline<>(name, this, this.calcTime, action);
  }

  // 可以并行执行的任务，各个的运行结果拼接到createResult创建的result中
  public <NewOut> Pipeline<Out, NewOut> thenMerge(String name,
                                                  Supplier<NewOut> createResult,
                                                  BiConsumer<String, Exception> onException,
                                                  MergeFunc<Out, NewOut>... mergeFuncs) {
    return new Pipeline<>(name, this, this.calcTime, (executorService, param, store) -> {
      WorkGroup workGroup = new WorkGroup(executorService);
      NewOut result = createResult.get();
      int pos = 0;
      for (MergeFunc<Out, NewOut> mergeFunc : mergeFuncs) {
        workGroup.add(name + ":" + pos, () -> mergeFunc.mergeExec(executorService, param, store, result));
      }
      workGroup.waitAllFinish(onException);
      return result;
    });
  }

  // 批量执行，返回一个数组
  public <Sub, NewOut> Pipeline<Out, List<NewOut>> thenBatch(String name,
                                                             BiConsumer<Out, Consumer<Sub>> looper,
                                                             BiConsumer<String, Exception> onException,
                                                             IPipeline<Sub, NewOut> action) {
    return new Pipeline<>(name, this, this.calcTime, (executorService, param, store) -> {
      WorkGroup workGroup = new WorkGroup(executorService);
      looper.accept(param, sub ->
          workGroup.add(name + ":" + sub, () -> action.execute(executorService, sub, store))
      );
      return workGroup.waitAllFinishAndGetResult(onException);
    });
  }

  public Out execute(ExecutorService executorService, IStore store) {
    return execute(executorService, 1, store);
  }

  @Override
  public Out execute(ExecutorService executorService, Object in, IStore store) throws RuntimeException {
    if (previous != null) {
      in = previous.execute(executorService, in, store);
    }

    // 如果上一步流程返回null表示流程终止，下面的都不再执行
    if (in == null) {
      return null;
    }

    long start = System.nanoTime();
    Out out = null;
    if (action != null) {
      out = action.execute(executorService, (In) in, store);
    }
    if (calcTime) {
      log.info("{} cost time = {} ms", name, (System.nanoTime() - start) / 1e6);
    }
    return out;
  }
}
