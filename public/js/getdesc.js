$(document).ready(function() {
    var md = window.markdownit();

    var descWrap = $('#task_desc .value');
    var href = location.href.split('/');
    var id = href[5];
    var type = href[3];

    $.get("/getdesc", { id: id, type:type }).done(function(data){
        var html = md.render(data);
        descWrap.html(html);
    })
});