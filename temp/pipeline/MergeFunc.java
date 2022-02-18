package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;

public interface MergeFunc<In, Out> {
  void mergeExec(ExecutorService executorService, In in, IStore store, Out out);
}
