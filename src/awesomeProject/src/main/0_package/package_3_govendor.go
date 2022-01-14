package main

/*
	1.用govendor的好处:
		1>.自 1.5 版本开始引入 govendor 工具，该工具, 通过读取该目录下的vendor.json(类似pom文件) 文件来记录依赖包的版本,
			将项目中, 所有依赖的外部包, 放到项目下的 vendor 目录下（对比 nodejs 的 node_modules 目录).
			这样可以方便用户使用相对稳定的依赖。
			所以,不同的人, 在不同的时刻,在同一个项目中,通过govendor下载的依赖,都会是同一个版本,
			而不会出现每个人下载的依赖的版本不同的情况.

  	2.关键:
		1>.项目下的vendor目录. 这个是固定的.
		2>.vendor/目录下的vendor.json文件. 里面记录了所依赖的包, 和包的版本.
		3>.使用: 在项目根目录下: #govendor sync 即可自动下载依赖. 根据的是上面的vendor.json文件.

    3.注意.
        用govendor下的包,是局部的, 只存在于其对应的项目中的,并不是全局系统级别的.
        即:
        对于 govendor 来说，主要存在三种位置的包：
        1>.项目自身的包组织为本地（local）包；
        2>.传统的存放在 $GOPATH 下的依赖包为外部（external）依赖包;(全局)
        3>.被 govendor 管理的放在 vendor 目录下的依赖包则为 vendor 包.(局部)


	4.常用命令
		1>.先安装govendor:
			#go get -u github.com/kardianos/govendor
		2>.将生成一个vendor目录, 和里面的vendor.json, 里面并没有依赖包信息.
			#govendor init
		3>.和pakeage.json同样, 我们需要将这个vendor.json添加到GIT, 但忽略vendor下的其他文件:
			.gitignore

			/vendor/*
			!/vendor/vendor.json

		4>.拉取依赖到vendor
			其他人可以使用vendor.json重新安装依赖包到vendor
			#govendor sync
 */