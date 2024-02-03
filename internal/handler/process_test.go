package handler

import (
	"context"
	"testing"
	"time"
	apiv3 "yyyoichi/Collo-API/internal/api/v3"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCoMatrixes(t *testing.T) {
	l, _ := time.LoadLocation("Asia/Tokyo")
	var initV3ReqConfig = func(v3req *apiv3.RequestConfig) *apiv3.RequestConfig {
		v3req.Keyword = "科学"
		v3req.From = timestamppb.New(time.Date(2023, 11, 1, 0, 0, 0, 0, l))
		v3req.Until = timestamppb.New(time.Date(2023, 11, 5, 0, 0, 0, 0, l))
		v3req.UseNdlCache = true
		v3req.CreateNdlCache = true
		return v3req
	}
	var testProcess = func(v3req *apiv3.RequestConfig) {
		config := NewConfig(initV3ReqConfig(v3req))
		handleErr := func(err error) {
			require.NoError(t, err)
		}
		var previousProcess float32
		handleProcessResp := func(process float32) {
			require.True(t, previousProcess < process)
			previousProcess = process
		}
		coMatrixes := NewCoMatrixes(
			context.Background(),
			ProcessHandler{
				Err:  handleErr,
				Resp: handleProcessResp,
			},
			config,
		)
		require.True(t, len(coMatrixes.Data) > 0)
		require.EqualValues(t, 1, previousProcess)
	}
	t.Run("GroupID is empty", func(t *testing.T) {
		t.Parallel()
		v3req := &apiv3.RequestConfig{
			ForcusGroupId: "",
		}
		testProcess(v3req)
	})
	t.Run("GroupID is 'total'", func(t *testing.T) {
		t.Parallel()
		v3req := &apiv3.RequestConfig{
			ForcusGroupId: "total",
		}
		testProcess(v3req)
	})
	t.Run("API is 'Meeting'", func(t *testing.T) {
		t.Parallel()
		v3req := &apiv3.RequestConfig{
			ForcusGroupId: "total",
			NdlApiType:    apiv3.RequestConfig_NDL_API_TYPE_MEETING,
		}
		testProcess(v3req)
	})
	t.Run("PickGroupID is 'Month'", func(t *testing.T) {
		t.Parallel()
		v3req := &apiv3.RequestConfig{
			PickGroupType: apiv3.RequestConfig_PICK_GROUP_TYPE_MONTH,
		}
		testProcess(v3req)
	})
	t.Run("Not found GroupID", func(t *testing.T) {
		t.Parallel()
		v3req := &apiv3.RequestConfig{
			ForcusGroupId: "not found forcus group id",
		}
		config := NewConfig(initV3ReqConfig(v3req))
		handleErr := func(err error) {
			require.NoError(t, err)
		}
		coMatrixes := NewCoMatrixes(
			context.Background(),
			ProcessHandler{
				Err:  handleErr,
				Resp: func(f float32) {},
			},
			config,
		)
		require.True(t, len(coMatrixes.Data) == 0)
	})
}
