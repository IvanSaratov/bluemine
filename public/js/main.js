function customStat(status) {
    switch (status) {
        case "Новая": {
            return 0;
            break;
        }
        case "В процессе": {
            return 1;
            break;
        }
        case "Отложена": {
            return 2;
            break;
        }
        case "Закрыта": {
            return 3;
            break;
        }
    }
}

function customPriority(priority) {
    switch (priority) {
        case "Низкий": {
            return 0;
            break;
        }
        case "Средний": {
            return 1;
            break;
        }
        case "Высокий": {
            return 2;
            break;
        }
    }
}

function MDParse(text, wrap) {
    var md = window.markdownit();
    var html = md.render(text);
    $(wrap).html(html);
}

function getTaskDescOrWikiArticle(type, id) {
    switch (type) {
        case 'task': {
            var descWrap = $('#task_desc .value');
            $.get("/get/taskdesc", { id: id }).done(function(data){
                MDParse(data, descWrap)
            })
            break
        }
        case 'wiki': {
            var artWrap = $('#wiki_desc .value');
            $.get("/get/wikiarticle", { id: id }).done(function(data){
                MDParse(data, artWrap)
            })
            break
        }
    }
}

function getWikiList(handle) {
    $.get("/get/wikilist", function(data) {
        handle(data);
    })
}

function fatherCount(data, obj) {
    var f = 0;
    if (obj.WikiFatherIDStr == "0") {
        return 0
    } else {
        for (var i = 0; i < data.length; i++) {
            if (data[i].WikiIDStr == obj.WikiFatherIDStr) {
                f += 1;
                return f += fatherCount(data, data[i]);
            }
        }
    }
}

function nodeExists(id) {
    var node = document.getElementById(id)
    return node != null;
}