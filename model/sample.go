package model

type SampleModel struct{}

var Sample SampleModel

func (t *SampleModel) GetList() (r string) {
	r = "Hello Go ~"
	return
}
