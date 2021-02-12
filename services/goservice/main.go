package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gopkg.in/olahol/melody.v1"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)

	r := gin.Default()

	m := melody.New()

	// curl -X POST 'localhost:5000/flespi?deviceToken=9ebbd0b25760557393a43064a92bae539d962103&deviceName=fake-device-911' -d "{"eventType": "AAS_PORTAL_START", "data": {"uid": "hfe3hf45huf33545", "aid": "1", "vid": "1"}}"
	r.POST("/flespi", func(c *gin.Context) {
		buf := make([]byte, 1024)
		num, _ := c.Request.Body.Read(buf)
		reqBody := string(buf[0:num])
		c.JSON(200, reqBody)
		println("\n[CUSTOM LOG]post /flespi req \n", reqBody)
	})

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/channel/:name", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "chan.html")
	})

	r.GET("/channel/:name/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.BroadcastFilter(msg, func(q *melody.Session) bool {
			return q.Request.URL.Path == s.Request.URL.Path
		})
	})

	r.GET("/user", func(c *gin.Context) {
		c.String(200, "popos")
	})

	v1user := r.Group("/user")
	{
		v1user.GET("/:name", func(c *gin.Context) {
			name := c.Param("name")
			c.String(http.StatusOK, "Hello %s", name)
		})

		// However, this one will match /user/john/ and also /user/john/send
		// If no other rs match /user/john, it will redirect to /user/john/
		v1user.GET("/:name/*action", func(c *gin.Context) {
			name := c.Param("name")
			action := c.Param("action")
			message := name + " is " + action
			c.String(http.StatusOK, message)
		})

		// For each matched request Context will hold the route definition
		v1user.POST("/:name/*action", func(c *gin.Context) {
			test := c.FullPath() == "/user/:name/*action" // true
			fmt.Print(test)
			c.JSON(200, test)
		})
	}

	decoup := r.Group("/decoup")
	{
		// /?webhookEventName=WEBHOOK_EVENT_NAME&deviceToken=DEVICE_TOKEN&deviceName=DEVICE_NAME
		decoup.POST("/:webhookEventName/:deviceToken/:deviceName", func(c *gin.Context) {
			// http.ServeFile(c.Writer, c.Request, "")
			webhookEventName := c.Param("webhookEventName")
			deviceToken := c.Param("deviceToken")
			deviceName := c.Param("deviceName")
			c.String(200, webhookEventName,deviceToken,deviceName)
		})
	}

	// r.GET("/flespi", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"deviceName": c.Query("deviceName"),
	// 		"deviceToken": c.Query("deviceToken"),
	// 	})
	// })
	go r.Run(":5000")
	go log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {

    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;

    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
    };

    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };

    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };

    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };

});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output"></div>
</td></tr></table>
</body>
</html>
`))
