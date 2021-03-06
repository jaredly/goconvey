package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"log"
)

type ConcurrentTester struct {
	shell     contract.Shell
	batchSize int
}

func (self *ConcurrentTester) SetBatchSize(batchSize int) {
	self.batchSize = batchSize
	log.Printf("Now configured to test %d packages concurrently.\n", self.batchSize)
}

func (self *ConcurrentTester) TestAll(folders []*contract.Package) {
	for _, folder := range folders {
		if folder.Active {
			self.shell.Execute("go", "test", "-i", folder.Name)
		}
	}

	if self.batchSize == 1 {
		self.executeSynchronously(folders)
	} else {
		newCuncurrentCoordinator(folders, self.batchSize, self.shell).ExecuteConcurrently()
	}
	return
}

func (self *ConcurrentTester) executeSynchronously(folders []*contract.Package) {
	for _, folder := range folders {
		if !folder.Active {
			log.Printf("Skipping execution: %s\n", folder.Name)
			continue
		}
		log.Printf("Executing tests: %s\n", folder.Name)
		folder.Output, folder.Error = self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder.Name)
		if folder.Error != nil && folder.Output == "" {
			panic(folder.Error)
		}
	}
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := &ConcurrentTester{}
	self.shell = shell
	self.batchSize = defaultBatchSize
	return self
}

const defaultBatchSize = 10
