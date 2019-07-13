$(document).ready(function() {
    var descWrap = $('#task_desc .value');
    var href = location.href.split('/');
    var id = href[5];
    $.get("/gettaskdesc", { id: id }).done(function(data){
        data = data.replace(/(?:\r\n|\r|\n)/g, '<br>');
        descWrap.html(data);
    })
});