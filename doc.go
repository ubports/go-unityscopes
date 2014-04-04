/*
Package scopes is used to write Unity scopes in Go.

Scopes are implemented through types that conform to the Scope interface.

    type MyScope struct {}

The shell may ask the scope for search results, which will cause the
Search method to be invoked:

    func (s *MyScope) Search(query string, reply *SearchReply cancelled <-chan bool) error {
        category := reply.RegisterCategory("cat_id", "category", "", "")
        result := NewCategorisedResult(category)
        result.SetTitle("Title")
        reply.Push(result)
        return nil
    }

In general, scopes will:

* Register result categories via reply.RegisterCategory()

* Create new results via NewCategorisedResult(), and push them with reply.Push(result)

* Check for cancellation requests via the provided channel.

The Search method will be invoked with an empty query when sufacing
results are wanted.

The shell may ask the scope to provide a preview of a result, which causes the Preview method to be invoked:

    func (s *MyScope) Preview(result *Result, reply *PreviewReply, cancelled <-chan bool) error {
        widget := NewPreviewWidget("foo", "text")
        widget.AddAttributeValue("text", "Hello")
        reply.PushWidgets(widget)
        return nil
    }

The scope should push one or more slices of PreviewWidgets using reply.PushWidgets.  PreviewWidgets can be created with NewPreviewWidget.

Additional data for the preview can be pushed with reply.PushAttr.

Finally, the scope can be exported in the main function:

    func main() {
        scopes.Run("scope-name", &MyScope{})
    }

The scope executable can be deployed to a scope directory named like:

    /usr/lib/${arch}/unity-scopes/${scope_name}

In addition to the scope executable, a scope configuration file named
${scope_name}.ini should be placed in the directory.  Its contents
should look something like:

    [ScopeConfig]
    DisplayName = Short name for the scope
    Description = Long description of scope
    Author =
    ScopeRunner = ${scope_executable}
*/
package scopes
