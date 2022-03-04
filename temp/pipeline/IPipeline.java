package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;
import java.util.function.Consumer;

public interface IPipeline<In, Out> {
  // store是个流程共享的对象，用于传递混存内容等
  void execute(ExecutorService executorService, In in, IStore store, Consumer<Out> callback) throws RuntimeException;
}
