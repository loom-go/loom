package rendering

// import "github.com/AnatoleLucet/loom"

// func RenderChildren(ctx *loom.NodeContext, parent any, children ...loom.Node) error {
// 	for i, child := range children {
// 		childCtx := ctx.Child(i)
// 		childCtx.SetParent(parent)
//
// 		if err := child.Render(childCtx); err != nil {
// 			return err
// 		}
// 	}
//
// 	// Trim any extra child contexts from previous renders
// 	ctx.TrimChildren(len(children))
//
// 	return nil
// }
