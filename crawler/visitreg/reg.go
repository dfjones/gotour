package visitreg

type VisitRegister interface {
  Visit(string)
  IsVisited(string) bool
  Close() map[string]bool
}
