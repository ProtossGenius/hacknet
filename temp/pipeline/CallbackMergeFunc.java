package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;

public interface CallbackMergeFunc<In, Out> {
  void mergeExec(ExecutorService executorService, In in, IStore store, Out out, CallbackGroup callbackGroup);
}
