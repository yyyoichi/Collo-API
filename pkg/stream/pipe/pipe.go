package pipe

import "context"

func Generator[T interface{}](cxt context.Context, seeds ...T) <-chan T {
	ch := make(chan T, len(seeds))
	go func() {
		defer close(ch)
		for _, s := range seeds {
			select {
			case <-cxt.Done():
				return
			case ch <- s:
			}
		}
	}()
	return ch
}

func GeneratorWithFn[I interface{}, O interface{}](cxt context.Context, fn func(I) O, seeds ...I) <-chan O {
	ch := make(chan O, len(seeds))
	go func() {
		defer close(ch)
		for _, s := range seeds {
			select {
			case <-cxt.Done():
				return
			case ch <- fn(s):
			}
		}
	}()
	return ch
}

type ChunkFnResp[O interface{}] struct {
	Out O   // アウトプット
	Len int // 探査数
}

// データセット[seeds] []I をチャンクごとに chan O にして返す。[fn]引数データセットから最初のチャンクと探査数を返す関数。
func Chunk[I interface{}, O interface{}](cxt context.Context, fn func([]I) ChunkFnResp[O], seeds ...I) <-chan O {
	outCh := make(chan O)
	fnRespCh := make(chan ChunkFnResp[O])
	go func() {
		defer close(outCh)
		// 探査位置
		i := 0
		for {
			target := seeds[i:]
			if len(target) == 0 {
				break
			}
			select {
			case <-cxt.Done():
			case fnRespCh <- fn(target):
				resp := <-fnRespCh
				outCh <- resp.Out
				i += resp.Len
			}
		}
	}()
	return outCh
}

func Line[I interface{}, O interface{}](cxt context.Context, inCh <-chan I, fn func(I) O) <-chan O {
	outCh := make(chan O)
	go func() {
		defer close(outCh)
		for {
			select {
			case <-cxt.Done():
				return
			case in, ok := <-inCh:
				if !ok {
					return
				}
				select {
				case <-cxt.Done():
					return
				case outCh <- fn(in):
				}
			}
		}
	}()
	return outCh
}
