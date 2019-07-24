getWikiList(function(data) {
    var wrap = $('#input_wiki_father');
    var html = '';
    for (var i = 0; i < data.length; i++) {
        if (data[i].WikiFatherID == 0) {
            html = '<option id="' + data[i].WikiIDStr + '">' + data[i].WikiName + '</option>'
            $(wrap).append(html);
        } else {
            html = '<option id="' + data[i].WikiIDStr + '"></option>'
            $(wrap).append(html);
            $('#' + data[i].WikiIDStr).insertAfter('#' + data[i].WikiFatherIDStr)
            var f = fatherCount(data, data[i])
            var str = '';
            for (var j = 0; j < f; j++) {
                str += '>>';
            }
            $('#' + data[i].WikiIDStr).text(str+data[i].WikiName)
        }
    }
})