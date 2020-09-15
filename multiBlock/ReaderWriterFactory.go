package multiBlock

import (
	log2 "github.com/bhbosman/gocommon/log"
	"go.uber.org/fx"
)

type IReaderWriterFactoryService interface {
	Create() iMultiBlock
	CreateAndAddBuffer(buffer []byte) (iMultiBlock, error)
}

type ReaderWriterFactoryService struct {
}

func (r ReaderWriterFactoryService) Create() iMultiBlock {
	return NewReaderWriter()
}

func (self *ReaderWriterFactoryService) CreateAndAddBuffer(buffer []byte) (iMultiBlock, error) {
	result := NewReaderWriterSize(len(buffer))
	_, err := result.Write(buffer)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func NewReaderWriterFactoryService() *ReaderWriterFactoryService {
	return &ReaderWriterFactoryService{}
}

type ReaderWriterFactoryInParams struct {
	fx.In
	LogFactory *log2.LogFactory
}
type ReaderWriterFactoryOutParams struct {
	fx.Out
	ReaderWriterFactory         IReaderWriterFactoryService
	ReaderWriterFactoryInstance *ReaderWriterFactoryService
}

func ProvideReaderWriterFactoryService() fx.Option {
	return fx.Provide(
		func(params ReaderWriterFactoryInParams) ReaderWriterFactoryOutParams {
			result := NewReaderWriterFactoryService()
			return ReaderWriterFactoryOutParams{
				ReaderWriterFactory:         result,
				ReaderWriterFactoryInstance: result,
			}
		})
}
