$(document).ready(function() {
    var descWrap = $('#task_desc .value');
    var href = location.href.split('/');
    var id = href[5];
    $.get("/gettaskdesc", { id: id }).done(function(data){
        descWrap.text(data);
    })
});