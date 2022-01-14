# tmm
快速搭建WEB站点以及提供RESTful接口 。

## 一：静态资源站点

	

``` go
	router := gin.Default()

	// 静态资源加载，本例为css,js以及资源图片
	router.StaticFS("/public", http.Dir("D:/goproject/src/github.com/ffhelicopter/tmm/website/static"))
	router.StaticFile("/favicon.ico", "./resources/favicon.ico")
```

## 二：动态站点
模板可调用静态资源站点的css，图片等

	

``` go
// 导入所有模板，多级目录结构需要这样写
	router.LoadHTMLGlob("website/tpl/*/*")
	
	// website分组
	v := router.Group("/")
	{

		v.GET("/index.html", handler.IndexHandler)
		v.GET("/add.html", handler.AddHandler)
		v.POST("/postme.html", handler.PostmeHandler)
	}
```

## 三：中间件的使用，在API中可能使用限流，身份验证等


	

``` go
	// 中间件 Go语言的net/http包特别容易构建中间件。
	// Gin提供了类似的中间件。需要注意的是中间件只对注册过的路由起作用。
	// 可以限定中间件的作用范围。大致分为全局中间件，单个处理程序中间件和组中间件。

	// 使用全局CORS中间件。
	// router.Use(Cors())
	// 即使是全局中间件，在use前的代码不受影响
	// 也可在handler中局部使用，见api.GetUser

	// 身份认证中间件，对于API，我们可以考虑JSON web tokens

	//rate-limit 限流中间件 
	lmt := tollbooth.NewLimiter(1, nil)
	//lmt := tollbooth.NewLimiter(1, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Hour})
	//lmt.SetHeader("X-Access-Token", []string{"abc123", "xyz098"})
	//lmt.SetBasicAuthUsers([]string{"sansa"})
	lmt.SetMessage("服务繁忙，请稍后再试...")
	//tollbooth.LimitByKeys(lmt, []string{"127.0.0.1", "/"})
```


## 四：RESTful API接口

``` go
// 组路由以及版本控制
	v1 := router.Group("/v1")
	{
		// 下面是组中间件的用法
		// v1.Use(Cors())

		// 单个中间件的用法
		// v1.GET("/user/:id/*action",Cors(), api.GetUser)

		// rate-limit
		v1.GET("/user/:id/*action", LimitHandler(lmt), api.GetUser)

		//v1.GET("/user/:id/*action", Cors(), api.GetUser)
		// AJAX OPTIONS ，下面是有关OPTIONS用法的示例
		// v1.OPTIONS("/users", OptionsUser)      // POST
		// v1.OPTIONS("/users/:id", OptionsUser)  // PUT, DELETE
		/*
			// 对应的handler中增加处理
			func OptionsUser(c *gin.Context) {
			    c.Writer.Header().Set("Access-Control-Allow-Methods", "DELETE,POST,PUT")
			    c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			    c.Next()
				...
			}

		*/
	}
```


## 六: 关于后端endpoint接口.
```

/*

说明:
	一：取匹配符的参数值.
		v1.GET("/user/:id/*action", api.GetUser)
		冒号:加上一个参数名组成路由参数。可以使用c.Params的方法读取其值。注意请求末尾不能加"/"，不能匹配
		除了:，gin还提供了*号处理参数，*号能匹配的规则就更多。
		"/user/:id/*action" 这里最基本的访问是/user/111/。注意末尾有"/"  //fixme 注意末尾有/与无/的两种形式.
		取*后参数 action := c.Param("action")

	二：***Query String***
		1.介绍:
			web提供的服务通常是client和server的交互。其中客户端向服务器发送请求，除了路由参数，
			其他的参数无非两种，查询字符串query string和报文体body参数。
			所谓query string，即路由用，用?以后连接的key1=value2&key2=value2的形式的参数。
			当然这个key-value是经过urlencode编码。

		2：取Query String的值
			query string,对于参数的处理，经常会出现参数不存在的情况，提供默认值处理：
			firstname := c.DefaultQuery("firstname", "Guest")
			lastname := c.Query("lastname")
			注意默认值是指参数没有出现在url中，如果出现但没有值，则为空字符串。
			例如:
				localhost:8080/   c.Query("id")值为空
				localhost:8080/?id=1   c.Query("id")值为1


	三：取form表单的值.
		0.介绍
			body http的报文体传输数据就比query string稍微复杂一点，常见的格式就有四种。
			例如:
				application/json，
				application/x-www-form-urlencoded,
				application/xml,
				multipart/form-data。
			后面一个主要用于图片上传。json格式的很好理解，urlencode其实也不难，无非就是把query string的内容，
			放到了body体里，同样也需要urlencode。默认情况下，c.PostFROM解析的是x-www-form-urlencoded或from-data的参数。

		1.取form中的字符串:
			html中的标签是这样的(字符串):
				<input type="checkbox" value="girl" name="message"/>
				<input type="checkbox" value="girl" name="nick"/>
			后端取值:
				message := c.PostForm("message")
				nick := c.DefaultPostForm("nick", "anonymous")
				post方法也提供了处理默认参数的情况。同理，如果参数不存在，将会得到空字串。

		2.取form中的单个文件:
			html中的传文件:
				<form action="http://localhost:8080/upload" method="post" enctype="multipart/form-data">
					头像:
					<input type="file" name="file">
					<br>
					<input type="submit" value="提交">
				</form>
			后端取文件:
				// single file
				file, _ := c.FormFile("file")
				log.Println(file.Filename)

				// Upload the file to specific dst.
				// func (c *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {}
				c.SaveUploadedFile(file, file.Filename)

				或:也可以直接使用io操作，拷贝文件数据。
				//out, err := os.Create(filename)
				//defer out.Close()
				//_, err = io.Copy(out, file)

				c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))

		3.取form中的多个文件:
            重点是html中的两处:
                1.enctype="multipart/form-data"
                2. <input .... multiple>
			html中的传文件:
				<!DOCTYPE html>
                <html lang="en">
                    <head>
                        <meta charset="UTF-8">
                        <title>文件s</title>
                    </head>
                    <body>
                    <h1>上传多个文件</h1>
                    
                    <form action="http://127.0.0.1:8080/upload" method="post" enctype="multipart/form-data">
                        Files: <input type="file" name="files" multiple><br><br>
                        <input type="submit" value="提交">
                    </form>
                    </body>
                </html>
			后端取文件:
				//所谓多个文件，无非就是多一次遍历文件，然后一次copy数据存储即可。
				func main() {
					router := gin.Default()
					// Set a lower memory limit for multipart forms (default is 32 MiB)
					router.MaxMultipartMemory = 8 << 20 // 8 MiB
					//router.Static("/", "./public")
					router.POST("/upload", func(c *gin.Context) {

						// Multipart form
						form, err := c.MultipartForm()
						if err != nil {
							c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
							return
						}
						files := form.File["files"]

						for _, file := range files {
							if err := c.SaveUploadedFile(file, file.Filename); err != nil {
								c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
								return
							}
						}

						c.String(http.StatusOK, fmt.Sprintf("Uploaded successfully %d files ", len(files)))
					})
					router.Run(":8080")
				}


	四:路由组Group

		0.介绍
			我们可以将拥有共同URL前缀的路由划分为一个路由组。
			习惯性一对{}包裹同组的路由，这只是为了看着清晰，你用不用{}包裹功能上没什么区别。

		1.普通的路由组
			func main() {
				r := gin.Default()
				userGroup := r.Group("/user")
				{
					userGroup.GET("/index", func(c *gin.Context) {...})
					userGroup.GET("/login", func(c *gin.Context) {...})
					userGroup.POST("/login", func(c *gin.Context) {...})

				}
				shopGroup := r.Group("/shop")
				{
					shopGroup.GET("/index", func(c *gin.Context) {...})
					shopGroup.GET("/cart", func(c *gin.Context) {...})
					shopGroup.POST("/checkout", func(c *gin.Context) {...})
				}
				r.Run()
			}

		2.普通的路由组
			路由组也是支持嵌套的，例如：
				shopGroup := r.Group("/shop")
					{
						shopGroup.GET("/index", func(c *gin.Context) {...})
						shopGroup.GET("/cart", func(c *gin.Context) {...})
						shopGroup.POST("/checkout", func(c *gin.Context) {...})
						// 嵌套路由组
						xx := shopGroup.Group("xx")
						xx.GET("/oo", func(c *gin.Context) {...})
					}

	五：后端返回json
    
        嵌套的json，只要嵌套gin.H就可以了。
    
           ********fixme 这是比较规范的返回值样式:*******
           c.JSON(http.StatusOK, gin.H{
                       "status":  gin.H{
                           "status_code": http.StatusOK,
                           "status":      "ok",
                       },
                       "message": message,
                       "nick":    nick,
                   })

	六：
        发送数据给服务端，并不是post方法才行，put方法一样也可以。
        同时querystring和body也不是分开的，两个同时发送也可以。
            router.PUT("/post", func(c *gin.Context) {
                id := c.Query("id")
                page := c.DefaultQuery("page", "0")
                name := c.PostForm("name")
                message := c.PostForm("message")
                fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)
                c.JSON(http.StatusOK, gin.H{
                    "status_code": http.StatusOK,
                })
            })


	七：参数绑定(请求参数-->VO)
	
        0.介绍:
            模型绑定可以将请求体绑定给一个类型。目前Gin支持JSON、XML、YAML和标准表单值的绑定。
            简单来说, 就是根据Body数据类型，将数据赋值到指定的结构体变量中 (类似于序列化和反序列化) 。

            Gin提供了两套绑定方法：
            
            Must bind
                方法：Bind,BindJSON,BindXML,BindQuery,BindYAML
                行为：这些方法使用MustBindWith。如果存在绑定错误，则用c终止请求，
                    使用c.AbortWithError (400) .SetType (ErrorTypeBind)即可。
                    将响应状态代码设置为400，Content-Type header设置为text/plain;charset = utf - 8。
                    请注意，如果在此之后设置响应代码，将会受到警告：[GIN-debug][WARNING] Headers were already written. Wanted to override status code 400 with 422将导致已经编写了警告[GIN-debug][warning]标头。
                    如果想更好地控制行为，可以考虑使用ShouldBind等效方法。
            Should bind
                方法：ShouldBind,ShouldBindJSON,ShouldBindXML,ShouldBindQuery,ShouldBindYAML
                行为：这些方法使用ShouldBindWith。如果存在绑定错误，则返回错误，开发人员有责任适当地处理请求和错误。
                注意，使用绑定方法时，Gin 会根据请求头中 Content-Type 来自动判断需要解析的类型。
                    如果你明确绑定的类型，你可以不用自动推断，而用 BindWith 方法。 
                    你也可以指定某字段是必需的。如果一个字段被 binding:"required" 修饰而值却是空的，请求会失败并返回错误。

            关于常用的:content-type。
                我们已经见识了x-www-form-urlencoded类型的参数处理，
                但是,现在越来越多的应用习惯使用JSON来通信,即: application/json.
                也就是无论返回的response还是提交的request，其content-type类型都是"application/json"的格式。
                但是, 对于一些旧的web表单页还是x-www-form-urlencoded的形式，这就也需要我们的服务器能hold住这多种content-type的参数了。
    
        1.关于用VO接参:

            1.VO字段设置非空(binding):
                针对下面VO(controller中接请求参数用的)的字段,可以看出:
                    1>.在结构体中，如果设置了binding标签的字段（如下面的username和password），
                        如果客户端没传就会抛错误(todo 待妥善处理这种错误, 看GO是否有全局exception handler这种东西)。
                    2>.对于非banding的字段（age），对于客户端没有传，User结构会用零值填充。
                        对于User结构没有的参数，会自动被忽略。
                这里结构体如内嵌将不能正常bind，***一句话要bind就不要使用内嵌***
            
            2.把请求参数绑定到VO中,用c.Bind(&VO):
    
               package main
            
                //注意结构体中的某些字段被binding了.
               type User struct {
                    //binding了, 不传过来值就会报错.
                    Username string `form:"username" json:"username" binding:"required"`
                    //binding了, 不传过来值就会报错.  
                    Passwd   string `form:"passwd" json:"passwd" bdinding:"required"`  
                    Age      int    `form:"age" json:"age"`
                   }
        
               func main() {
                router := gin.Default()
                router.POST("/login", func(c *gin.Context) {
                    var user User
                    var err error
                    contentType := c.Request.Header.Get("Content-Type")
        
                    //处理了两种contentType.(fixme:下面有更高级的用法. Gin内置支持的.)
                    switch contentType {
                        case "application/json":
                            //把参数绑定到VO里, 是用API:c.BindJSON(&VO)这么做的, 不是放参数上自动的.
                            //下面有更高级的用法: c.Bind()
                            err = c.BindJSON(&user)
                        case "application/x-www-form-urlencoded":
                            err = c.BindWith(&user, binding.Form)
                    }
        
                    if err != nil {
                        fmt.Println(err)
                        log.Fatal(err)
                    }
        
                    //不推荐这样直接返回数据, 要返回符合规范的数据. 带message和code的. 上面有提到.
                    c.JSON(http.StatusOK, gin.H{
                        "user":   user.Username,
                        "passwd": user.Passwd,
                        "age":    user.Age,
                    })
                })
               }
        
                //当然，gin还提供了更加高级方法，c.Bind，它会更加content-type自动推断是bind表单还是json的参数。
                //fixme: 用这种高级方式比较好. 不要用上面switch版本的.
                router.POST("/login", func(c *gin.Context) {
                    var user User
                    err := c.Bind(&user)
                    if err != nil {
                        fmt.Println(err)
                        log.Fatal(err)
                    }
        
                    //不推荐这样直接返回数据, 要返回符合规范的数据. 带message和code的. 上面有提到.
                    c.JSON(http.StatusOK, gin.H{
                        "username": user.Username,
                        "passwd":   user.Passwd, "age": user.Age,
                    })
                })
        
            // c.ShouldBind c.ShouldBindBodyWith 这两个消耗c.Request.Body。 他们也可以使用，但不能多次调用。
            // 由于c.ShouldBindBodyWith在bind前会存储body到context，性能有影响，
            // 但是只对于JSON, XML, MsgPack, ProtoBuf格式。
            // 对于:Query, Form, FormPost, FormMultipart 等格式，可以多次调用，性能也不受影响。
        
                if errA := c.ShouldBindBodyWith(&objA, binding.JSON); errA == nil {
                  c.String(http.StatusOK, `the body should be formA`)
                  // At this time, it reuses body stored in the context.
                } else if errB := c.ShouldBindBodyWith(&objB, binding.JSON); errB == nil {
                  c.String(http.StatusOK, `the body should be formB JSON`)
                  // And it can accepts other formats
                } else if errB2 := c.ShouldBindBodyWith(&objB, binding.XML); errB2 == nil {
                  c.String(http.StatusOK, `the body should be formB XML`)
                } else {
                  ...
                }


	九：服务端返回值格式的处理.
        多格式渲染
        既然请求可以使用不同的content-type，响应也如此。
        通常响应会有: html，text，plain，json和xml等。
        gin提供了很优雅的渲染方法。
        主要有: c.String， c.JSON，c.HTML，c.XML。 (常用)


	十：重定向 
        1.重定向到其他url:
	        c.Redirect(http.StatusMovedPermanently, "https://google.com")
        2.重定向到其他endpoint(路由):
            //对请求的URI进行修改
            c.Request.URL.Path = "/b"
            //继续后续的处理
            r.HandleContext(c)


	十一：异步
        异步 gin里可以借助协程实现异步任务。fixme: 注意: 因为涉及异步过程，
        请求的上下文需要copy到异步的上下文，并且这个上下文是只读的。
        某些情况下非常有用，后面会看到这方面的应用
            cCp := c.Copy()
            go func() {
                time.Sleep(5 * time.Second)
                log.Println("Done! in path" + cCp.Request.URL.Path)
            }()

*/

```

## 七: 个人进度.

``` go

TODO:
    1.CRUD API --DONE. 见:src/com.jsflzhong/6_web/Gin/website_1/api/restAPI.go
    2.REST API ALL Ongoing. 见:src/com.jsflzhong/6_web/Gin/website_1/api/restAPI.go
    3.Exception handler
        可以考虑结合上面的VO中的bingding.
    4.Load Conf
    5.Transaction?    
    6.Log File?
    
```

## 九: Gin关键结构说明.

``` go

1.生成engine
2.加载静态资源目录.
3.注册路由器endpoints.
4.写路由器对应的hanlder.(相当于controller中的handler). 并把一些业务逻辑也写在这里.
5.把CRUD的SQL逻辑写在entity文件中.
6.把对DB的直接操作, 写在单独的go文件中. 被上面调用.
```


## 八: 目录说明

``` go

1.src/com.jsflzhong:
    存放主要的代码.
    包含一些basic功能的代码. 但是大部分basic功能的执行文件main不在这里. 放在了下面src/main下面.
    但是一些小型项目除外, 例如 src/com.jsflzhong/6_web/Gin/website_1 
        这个小型web项目的main文件在:src/com.jsflzhong/6_web/Gin/website_1/myMain.go    
    
2. src/main
    一些上面包中测试代码的main执行文件的位置.
    为了方便, 这里只放main执行文件.    

```
