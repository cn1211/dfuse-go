package dfuse

//
//import "github.com/mohae/deepcopy"
//
//type Callback struct {
//	processors []*CallbackProcessor
//}
//
//func (c *Callback) clone() *Callback {
//	return &Callback{
//		processors: deepcopy.Copy(c.processors).([]*CallbackProcessor),
//	}
//}
//
//type CallbackProcessor struct {
//	name    string
//	before  string
//	after   string
//	replace bool
//	remove  bool
//	//kind    string
//
//	processor func()
//	parent    *Callback
//}
//
//func (c *CallbackProcessor) Before(callbackName string) *CallbackProcessor {
//	c.before = callbackName
//	return c
//}
//
//func (c *CallbackProcessor) After(callbackName string) *CallbackProcessor {
//	c.after = callbackName
//	return c
//}
//
//func (c *CallbackProcessor) Register(callbackName string, callback func()) {
//	c.name = callbackName
//	c.processor = callback
//	c.parent.processors = append(c.parent.processors, c)
//}
//
//func (c *CallbackProcessor) Remove(callbackName string) {
//	c.name = callbackName
//	c.processor = func() {}
//	c.remove = true
//	c.parent.processors = append(c.parent.processors, c)
//}
//
//func (c *CallbackProcessor) Replace(callbackName string) {
//	c.name = callbackName
//	c.processor = func() {}
//	c.replace = true
//	c.parent.processors = append(c.parent.processors, c)
//}
//
//func (c *CallbackProcessor) Get(callbackName string) (callback func()) {
//	for _, process := range c.parent.processors {
//		if process.name == callbackName && !process.remove {
//			return process.processor
//		}
//	}
//	return nil
//}
