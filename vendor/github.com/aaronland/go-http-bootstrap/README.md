# go-http-bootstrap

Go HTTP middleware for Bootstrap (4.6.0)

## Important

* Documentation is incomplete.

* Because Bootstrap 4.x requires the use of jQuery for its JavaScript functionality this package does _not_ include or inject references to `bootstrap.js` by default. If you want to do so you will need to update the `BootstrapOptions.JS` array, like this:

```
	my_handler := ...	// valid http.HandlerFunc
	
	bootstrap_opts := bootstrap.DefaultBootstrapOptions()
	bootstrap_opts.JS = []string{"/javascript/bootstrap.min.js"}

	my_handler = bootstrap.AppendResourcesHandler(handler, bootstrap_opts)
```

## Example

See [example/main.go](example/main.go) for a currently-working example.

## See also

* https://getbootstrap.com/
