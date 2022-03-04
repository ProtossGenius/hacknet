package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;

public interface ThenFunc<In, Out> {
  Out execute(ExecutorService executorService, In in, IStore store) throws RuntimeException;
}
