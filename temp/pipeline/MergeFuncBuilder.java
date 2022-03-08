package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;
import java.util.function.BiConsumer;
import java.util.function.Consumer;

public class MergeFuncBuilder<In, Out> {
  private final CallbackGroup callbackGroup;
  private final ExecutorService executorService;
  private final In in;
  private final Out out;
  private final IStore store;

  public MergeFuncBuilder(ExecutorService service, In in, Out out, IStore store) {
    this.executorService = service;
    this.in = in;
    this.out = out;
    this.store = store;
    callbackGroup = new CallbackGroup(service);
  }

  public MergeFuncBuilder<In, Out> and(String name, MergeFunc<In, Out> mergeFunc) {
    callbackGroup.add(name, () -> mergeFunc.mergeExec(executorService, in, store, out));
    return this;
  }

  public MergeFuncBuilder<In, Out> and(String name, CallbackMergeFunc<In, Out> mergeFunc) {
    callbackGroup.add(name, () -> mergeFunc.mergeExec(executorService, in, store, out, callbackGroup), 2);
    return this;
  }

  public void execute(BiConsumer<String, Exception> onException, Consumer<Out> callback) {
    callbackGroup.execute(onException, nth -> callback.accept(out));
  }
}
