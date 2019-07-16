$(document).ready(function() {
    var descWrap = $('#task_desc .value');
    var href = location.href.split('/');
    var id = href[5];
    var type = href[3];

    $.get("/getdesc", { id: id, type:type }).done(function(data){
        data = data.replace(/(?:\r\n|\r|\n)/g, '<br>');
        descWrap.html(data);
    })
});