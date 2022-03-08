package com.yqg.recall.common.util.pipeline;

import lombok.extern.slf4j.Slf4j;

import java.util.List;
import java.util.concurrent.ExecutorService;
import java.util.function.BiConsumer;
import java.util.function.Consumer;
import java.util.function.Function;
import java.util.function.Supplier;

/**
 * 流程（管道）
 *
 * @param <In>  入参
 * @param <Out> 出参
 */
@Slf4j
public class Pipeline<In, Out> implements IPipeline<In, Out> {
  private final static Object DEFAULT_PARAM = 1;
  // 流程名
  String name;
  // 流程要做的事
  IPipeline<In, Out> action;
  // 前一个流程
  Pipeline previous;
  // 是否计算执行时间
  boolean calcTime = false;

  /**
   * 用于第一个流程，第一个流程或许是不需要入参的
   *
   * @param name     流程名
   * @param calcTime 是否计算执行时间
   * @param action   行为
   */
  public Pipeline(String name, boolean calcTime, Function<ExecutorService, Out> action) {
    this.calcTime = calcTime;
    this.name = name;
    this.action = (executorService, in, store, cb) -> cb.accept(action.apply(executorService));
  }

  public Pipeline(String name, Function<ExecutorService, Out> action) {
    this(name, false, action);
  }

  /**
   * @param name     流程名
   * @param calcTime 是否计算执行时间
   * @param action   行为
   */
  public Pipeline(String name, boolean calcTime, ThenFunc<In, Out> action) {
    this.calcTime = calcTime;
    this.name = name;
    this.action = parseThenFunc(action);
  }

  public Pipeline(String name, ThenFunc<In, Out> action) {
    this(name, false, action);
  }

  /**
   * 用于连接流程
   *
   * @param name     流程名
   * @param previous 前一个流程
   * @param calcTime 是否计算执行时间
   * @param action   行为
   */
  private Pipeline(String name, Pipeline previous, boolean calcTime, IPipeline<In, Out> action) {
    this.name = name;
    this.previous = previous;
    this.calcTime = calcTime;
    this.action = action;
  }

  private static <Param, Result> IPipeline<Param, Result> parseThenFunc(ThenFunc<Param, Result> func) {
    return (executorService, param, store, callback) -> callback.accept(func.execute(executorService, param, store));
  }

  /**
   * 设置是否计算流程耗时
   *
   * @param calcTime
   */
  public void setCalcTime(boolean calcTime) {
    if (calcTime == this.calcTime) {
      return;
    }

    this.calcTime = calcTime;
    if (this.previous != null) {
      this.previous.setCalcTime(calcTime);
    }
  }

  /**
   * 进入下一个流程（在执行时，下一个流程的入参是当前流程的出参）
   *
   * @param name     下一个流程名
   * @param action   下一个流程的行为
   * @param <NewOut> 下一个流程的出参
   * @return 下一个流程
   */
  public <NewOut> Pipeline<Out, NewOut> then(String name, ThenFunc<Out, NewOut> action) {
    return new Pipeline<>(name, this, this.calcTime, parseThenFunc(action));
  }

  public <NewOut> Pipeline<Out, NewOut> then(String name, Pipeline action) {
    action.previous.name = name;
    action.previous = this;
    return action;
  }

  /**
   * 可以拆分执行的任务，各个的运行结果拼接到createResult创建的result中
   *
   * @param name         下一个流程名
   * @param createResult 因为是拼接到结果中，所以需要提前生成出参
   * @param onException  处理异常
   * @param mergeFuncs   被拆分的任务，注意避免多个函数处理result的同一个成员所导致的线程安全问题。
   * @param <NewOut>     下一个流程的出参
   * @return 下一个流程
   */
  @SafeVarargs
  final public <NewOut> Pipeline<Out, NewOut> thenMerge(String name,
                                                        Supplier<NewOut> createResult,
                                                        BiConsumer<String, Exception> onException,
                                                        MergeFunc<Out, NewOut>... mergeFuncs) {
    return new Pipeline<>(name, this, this.calcTime, (executorService, param, store, cb) -> {
      NewOut result = createResult.get();
      MergeFuncBuilder<Out, NewOut> mergeFuncBuilder = new MergeFuncBuilder<>(executorService, param, result, store);
      int pos = 0;

      for (MergeFunc<Out, NewOut> mergeFunc : mergeFuncs) {
        mergeFuncBuilder.and(name + ":" + (pos++), mergeFunc);
      }

      mergeFuncBuilder.execute(onException, cb);
    });
  }

  // 同上，只不过用了builder
  public <NewOut> Pipeline<Out, NewOut> thenMerge(String name,
                                                  Supplier<NewOut> createResult,
                                                  BiConsumer<String, Exception> onException,
                                                  Consumer<MergeFuncBuilder<Out, NewOut>> builder) {
    return new Pipeline<>(name, this, this.calcTime, (executorService, param, store, cb) -> {
      NewOut result = createResult.get();
      MergeFuncBuilder<Out, NewOut> mergeFuncBuilder = new MergeFuncBuilder<>(executorService, param, result, store);
      builder.accept(mergeFuncBuilder);
      mergeFuncBuilder.execute(onException, cb);
    });
  }

  /**
   * 下一个流程将批量处理当前流程的产出，并返回一个数组
   *
   * @param name        下一个流程名
   * @param looper      遍历函数，有两个入参，入参1是当前流程的产出，入参2是将入参1分解为批量处理的单元对象。
   *                    函数的任务是对入参1中每个需要批量处理的单元对象执行 Consumer<Sub>
   * @param onException 异常处理
   * @param action      对单元类型的处理
   * @param <Sub>       下个流程入参可被批量处理的单元类型
   * @param <NewOut>    下个流程的产出单元类型，下个流程的产出是该产出单元类型的List
   * @return 下个流程
   */
  public <Sub, NewOut> Pipeline<Out, List<NewOut>> thenBatch(String name,
                                                             BiConsumer<Out, Consumer<Sub>> looper,
                                                             BiConsumer<String, Exception> onException,
                                                             BatchFunc<Out, NewOut, Sub> action) {
    return new Pipeline<>(name, this, this.calcTime, (executorService, param, store, cb) -> {
      CallbackGroup callbackGroup = new CallbackGroup(executorService);
      looper.accept(param, sub ->
          callbackGroup.add(name + ":" + sub, () -> action.batchExec(executorService, param, sub, store))
      );

      callbackGroup.execute(onException, (List<NewOut> list) -> {
        cb.accept(list);
      });
    });
  }

  /**
   * 无入参执行流程（就算没有入参也需要传入一个非null的值）
   *
   * @param executorService 线程池
   * @param store           所有流程共用的数据类
   * @return 出参（注意，这是最后一个流程的出参）
   */
  public void execute(ExecutorService executorService, IStore store, Consumer<Out> callback) {
    execute(executorService, (In) DEFAULT_PARAM, store, callback);
  }

  public void execute(ExecutorService executorService, Object param, IStore store) {
    execute(executorService, (In) param, store);
  }

  public void execute(ExecutorService executorService, IStore store) {
    execute(executorService, store, nth -> {
    });
  }

  public Out blockExecute(ExecutorService executorService, Object in, IStore iStore) throws RuntimeException {
    final Box<Out> box = new Box<>();
    NonReentrantLockByWait lock = new NonReentrantLockByWait();
    lock.lock();
    execute(executorService, (In) in, iStore, out -> {
      lock.unlock();
      box.setValue(out);
    });
    lock.lock();
    lock.unlock();
    return box.getValue();
  }

  /**
   * 带参数执行流程
   * in 为null表示跳过该流程及所有后续流程
   *
   * @param executorService 线程池
   * @param in              入参（注意，这是第一个流程的入参）
   * @param store           所有流程共用的处理类
   * @return 出参（注意，这是最后一个流程的出参）
   * @throws RuntimeException 运行时错误，onException可能抛出运行时错误，用以打断流程
   */
  @Override
  public void execute(ExecutorService executorService, In in, IStore store, Consumer<Out> callback) throws
      RuntimeException {
    if (previous != null) {
      executorService.submit(() -> previous.execute(executorService, in, store, param ->
          doAction(executorService, (In) param, store, callback)
      ));

      return;
    }

    doAction(executorService, (In) in, store, callback);

  }

  private void doAction(ExecutorService executorService, In param, IStore store, Consumer<Out> callback) {
    executorService.submit(() -> {
      if (param == null) {
        callback.accept(null);
        return;
      }
      long start = System.nanoTime();

      try {
        action.execute(executorService, param, store, callback);
      } catch (Exception e) {
        log.error("in pipeline {}, error happened", name, e);
        callback.accept(null);
      }

      if (calcTime) {
        log.info("{} pipeline cost time = {} ms", name, (System.nanoTime() - start) / 1e6);
      }
    });
  }
}
