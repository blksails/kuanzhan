package kuanzhan

import "testing"

var testJson = `{
    "taskCreateTime": 1751984111632,
    "failedPages": [
        {
            "pageId": 1999095979,
            "status": "FAILED",
            "errorMsg": "当前页面已经删除",
            "siteId": 1394020996
        },
        {
            "pageId": 1913224245,
            "status": "FAILED",
            "errorMsg": "当前页面已经删除",
            "siteId": 8797234515
        },
        {
            "pageId": 3138030235,
            "status": "FAILED",
            "errorMsg": "当前页面已经删除",
            "siteId": 7644941491
        }
    ],
    "waitingPages": [],
    "succeedPages": [
        {
            "pageId": 2705120314,
            "status": "SUCCESS",
            "errorMsg": "成功",
            "siteId": 1394020996
        },
        {
            "pageId": 2636889128,
            "status": "SUCCESS",
            "errorMsg": "成功",
            "siteId": 8797234515
        },
        {
            "pageId": 1981238535,
            "status": "SUCCESS",
            "errorMsg": "成功",
            "siteId": 7644941491
        }
    ],
    "taskStatus": "PART_FAILED"
}`

func TestBatchModifyPagePublishPageJsData_UnmarshalJSON(t *testing.T) {
	type fields struct {
		TaskId string
		Task   Task
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test",
			fields: fields{
				TaskId: "123",
				Task:   Task{},
			},
			args: args{
				data: []byte(testJson),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &BatchModifyPagePublishPageJsData{
				TaskId: tt.fields.TaskId,
				Task:   tt.fields.Task,
			}
			if err := m.UnmarshalJSON(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("BatchModifyPagePublishPageJsData.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
