
23:27 <+kevlar> so you can set up your own listen routine and pass the conn 
          into ServeEchoService
23:28 <+kevlar> and you can use tls.Dial and then pass that conn to 
          NewEchoService (to replace DIalEchoService)

https://github.com/kylelemons/go-rpcgen

