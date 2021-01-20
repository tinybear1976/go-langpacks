# go-langpacks
**version: 0.9.1**



一个简单的多语言包系统

Simple multi language pack for golang。



# Install

```bash
go get  github.com/tinybear1976/go-langpacks
```

# Summary

1. Lang files
   - 默认语言包文件后缀 `.lps` (Lang Packs)
   - 语言包文件(.lps)开始第一行必须为语言标识，该标识配合您的程序作为运行时决定加载哪种语言包的关键内容，例如，我个人的习惯
   | 标识 | 含义 |
| ---- | ---- |
| zh   | 中文 |
| en   | 英文 |
   - 文件内(.lps)从第二行开始均为键值对，每组(行)键值对的分隔符为 `~`.  (e.g.)  `1000~show text`
   - 分隔符左边或右边允许存在空格，在加载后程序会自动去除空格内容.   (e.g.) `1000  ~  show text`
   - 不管如何书写键值对内容，程序在加载时总是期望：每行通过分隔符可以获得key（第一部分）和value（第二部分）两个部分，并且第一部分必须为整数(int)。如果程序在加载过程中，循环到某行并不能获得这个期望，程序将放弃该行内容（视为无效）
   
   lps文件演示：
   ```text
   en
   41501 ~ the user name is wrong or does not exist.
   41502 ~ wrong password.
   41503 ~ an error occurred while updating the token.
   41504 ~ The required post data was not submitted.
   ```
   
2. Loading mode

   ```go
   SetLoadMode(InMemory)    //or  SetLoadMode(InRedis)
   ```

   

   `InMemory`(默认) or `InRedis`  两种模式

   InMemory 模式: 适合少量的翻译文本内容，内部每个语言包都会采用一组 map [int] string来承载

   InRedis 模式: 适合大量的内容，并且多个程序共享语言包内容。如果采用该模式，则加载动作只执行一次即可，其他程序需要配置它们为该模式即可（除非将不同程序的语言包装载到不同的Redis上）

3. Load [Optional]

   - 装载需要给出语言包的基本路径（绝对路径,如果不给出路径则采用程序运行的当前目录），同时给出语言包的后缀名，以及内容分隔符号，这三部分内容如果全部传递空字符串，则表示采用默认值，即 路径=`当前程序运行路径`，后缀名=`.lps` ，分隔符=`~`
   - 加载后，会返回两个数值，一个表示预测文件中一共有多少词条需要加载，另外一个表示实际加载成功了多少个词条

4. Use

   通过前端程序传入对应的语言标识与具体文本id即可获得对应的文本，为了避免不必要的复杂操作，当没有检索到对应文本时，函数不会返回error，而是直接返回空字符串。用户可以根据自己的实际定义去检查语言包文件的键值对定义是否正确

5. 

# Example

```go
///   1. 初始化环境
InitLangPacksDefault()   //文件路径： ./  分隔符： ~   语言包文件后缀： .lps   加载模式：  InMemory
InitLangPacksDefaultRedis("127.0.0.1:6379", "password", 0) //文件路径： ./  分隔符： ~   语言包文件后缀： .lps   加载模式：  InRedis
InitLangPacks(lpsPath, lpsSuffix , separator , ipWithPort, pwd , db)

///   2. 设置模式[可选]
SetLoadMode(InRedis)   // 或  SetLoadMode(InMemory)

///   3. 加载语言包
Load()    // 返回执行结果 []LoadResult ，error，如果正常执行error==nil

///   4. 查询文本
Query("en", 1000)  // Query(langTag string, textId int) (str string)
```

# Tips

手工执行`InitLangPacksDefaultRedis` 或 `SetLoadMode` 后都将引发模式状态的变化，为了防止发生异常，需要手工调用一次`Load`。因此，为了防止运行时错误，当模式状态产生变化时，内部状态值`is_loaded=false`，在这种情况下，每次调用查询都将被强制驳回。


