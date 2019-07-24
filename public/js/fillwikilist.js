getWikiList(function(data) {
    var html = '';
    for (var i = 0; i < data.length; i++) {
        if (data[i].WikiFatherID == 0) {
            if (!nodeExists(data[i].WikiIDStr)) {
                var wrap = $('#wikilist #main');
                html = '<div id="' + data[i].WikiIDStr + '"><div class="item"><a href="/wiki/show/' + data[i].WikiIDStr + '">' + data[i].WikiName + '</a></div></div>'
                $(wrap).append(html);
                $('#' + data[i].WikiIDStr).children('.item').css("padding-left", '15px')
            }
        } else {
            if (!nodeExists(data[i].WikiIDStr)) {
                var wrap = $('#wikilist #main #' + data[i].WikiFatherIDStr);
                if (!nodeExists('child_of_' + data[i].WikiFatherIDStr + '')) {
                    $(wrap).append('<div class="nested" id="child_of_' + data[i].WikiFatherIDStr + '"></div>')
                }
                var wrap = $('#wikilist #main #' + data[i].WikiFatherIDStr + ' #child_of_' + data[i].WikiFatherIDStr);
                html = '<div id="' + data[i].WikiIDStr + '"><div class="nested_item"><a href="/wiki/show/' + data[i].WikiIDStr + '">' + data[i].WikiName + '</a></div></div>'
                $(wrap).append(html);
                var par = $('#' + data[i].WikiIDStr).parents().length
                var padding = 15 + 10*((par - 8)/2 + 1);
                $('#' + data[i].WikiIDStr).children('.nested_item').css("padding-left", padding.toString() + 'px')
            }
        }
    }
})