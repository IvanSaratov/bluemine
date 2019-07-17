$.get("/getwikilist", function(data) {
    var html = '';
    for (var i = 0; i < data.length; i++) {
        if (data[i].WikiFatherID == 0) {
            if (!nodeExists(data[i].WikiIDStr)) {
                var wrap = $('#wikilist #main');
                html = '<div id="' + data[i].WikiIDStr + '"><div class="item"><a href="/wiki/show/' + data[i].WikiIDStr + '">' + data[i].WikiName + '</a></div></div>'
                $(wrap).append(html);
            }
        } else {
            if (!nodeExists(data[i].WikiIDStr)) {
                var wrap = $('#wikilist #main #' + data[i].WikiFatherIDStr);
                $(wrap).append('<div class="nested" id="child_of_' + data[i].WikiFatherIDStr + '"></div>')
                var wrap = $('#wikilist #main #' + data[i].WikiFatherIDStr + ' #child_of_' + data[i].WikiFatherIDStr);
                html = '<div id="' + data[i].WikiIDStr + '"><div class="nested_item"><a href="/wiki/show/' + data[i].WikiIDStr + '">' + data[i].WikiName + '</a></div></div>'
                $(wrap).append(html);
            }
        }
    }
})

function nodeExists(id) {
    var node = document.getElementById(id)
    return node != null;
}