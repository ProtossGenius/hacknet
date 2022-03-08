package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;
import java.util.function.Consumer;

public interface CallbackBatchFunc<In, Out, Sub> {
  void batchExec(ExecutorService executorService, In in, Sub sub, IStore store, Consumer<Out> out) throws RuntimeException;
}
