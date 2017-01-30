sendMsg("filehelper", "JS脚本已经重新加载.")

var atMe = function (cb) {
    return function (group, user, content) {
        if (content.indexOf('@' + me.NickName) !== -1) {
            cb(group, user, content)
        }
    }
}

var onMessage = function (group, user, content) {
    if (atMe(content)) {
        if (contains(content, "reload")) {
            reloadJS()
            sendMsg(group, "好的")
            return
        }
        if (contains(content, "我是")) {
            reloadJS()
            sendMsg(group, "好的")
            return
        }
        sendMsg(group, "喊我干啥！")
        return
    }


    if (contains(content, "reload")) {
        reloadJS()
        sendMsg(group, "好的")
        return
    }
}

hear("reload", atMe(function (group, user, content) {
    reloadJS()
    sendMsg(group, "好的")
}))

hear("几点了", atMe(function (group, user, content) {
    sendMsg(group, new Date())
}))

hear("\\d+秒后告诉我.*", atMe(function (group, user, content) {
    var result = /(\d+)秒后告诉我(.*)/.exec(content)
    var after = result[1]
    var act = result[2]
    sendMsg(group, "好的！" + after + "秒后告诉你:" + act)
    setTimeout(function () {
        sendMsg(group, act)
    }, after * 1000)
}))

hear("db\\[.*\\]=.*", atMe(function (group, user, content) {
    var result = /db\[(.*)\]=(.*)/.exec(content)
    var dbName = result[1]
    var dbAddr = result[2]
    set("db:" + dbName, dbAddr)
    sendMsg(group, "好的！")
}))


hear("query\\[(.*)\\]:(.*)", atMe(function (group, user, content) {
    var result = /query\[(.*)\]:(.*)/.exec(content)
    var dbName = result[1]
    var sql = result[2]

    var dbAddr = get("db:" + dbName)

    sendMsg(group, "result:\n" + query(dbAddr, sql))
}))

hear("query\\[(.*)\\](.*)=(.*)", atMe(function (group, user, content) {
    var result = /query\[(.*)\](.*)=(.*)/.exec(content)
    var dbName = result[1]
    var sqlAlias = result[2]
    var sql = result[3]
    var dbAddr = get("db:" + dbName)
    set("sql:" + dbName + "/" + sqlAlias, sql)
    sendMsg(group, "result:\n" + query(dbAddr, sql))
}))


hear("query\\[(.*)\\](.*);", atMe(function (group, user, content) {
    var result = /query\[(.*)\](.*);/.exec(content)
    var dbName = result[1]
    var sqlAlias = result[2]
    var sql = get("sql:" + dbName + "/" + sqlAlias)
    var dbAddr = get("db:" + dbName)
    sendMsg(group, "result:\n" + query(dbAddr, sql))
}))