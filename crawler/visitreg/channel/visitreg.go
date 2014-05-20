package channel

func New() visitMap {
  vm := make(visitMap)
  go vm.run()
  return vm
}

type visitMap chan visitData

type visitData struct {
  action visitAction
  url string
  result chan<- bool
  data chan<- map[string]bool
}

type visitAction int

const (
  add visitAction = iota
  has
  end
)

func (vm visitMap) run() {
  store := make(map[string]bool)
  for command := range vm {
    switch command.action {
      case add:
        store[command.url] = true
      case has:
        _, found := store[command.url]
        command.result <- found
      case end:
        command.data <- store
        close(vm)
    }
  }
}

func (vm visitMap) Visit(url string) {
  vm <- visitData{action: add, url: url}
}

func (vm visitMap) IsVisited(url string) bool {
  reply := make(chan bool)
  vm <- visitData{action: has, url: url, result: reply}
  return <-reply
}

func (vm visitMap) Close() (map[string]bool){
  data := make (chan map[string]bool)
  vm <- visitData{action: end, data: data}
  return <-data
}
