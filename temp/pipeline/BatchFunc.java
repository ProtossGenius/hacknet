package com.yqg.recall.common.util.pipeline;

import java.util.concurrent.ExecutorService;

public interface BatchFunc<In, Out, Sub> {
  Out batchExec(ExecutorService executorService, In in, Sub sub, IStore store) throws RuntimeException;
}
