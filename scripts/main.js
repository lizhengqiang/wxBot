console.log(JSON.stringify(me))
sendMsg("filehelper", "JS脚本已经重新加载.")


var atMe = function (content) {
    return content.indexOf('@' + me.NickName) !== -1
}

var contains = function (str, s) {
    return str.indexOf(s) !== -1
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

hear("猫猫.*", onMessage)

hear(".*猫猫.*几点了", function (group, user, content) {
    sendMsg(group, new Date())
})