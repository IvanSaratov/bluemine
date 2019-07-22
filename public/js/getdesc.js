$(document).ready(function() {
    var md = window.markdownit();

    var descWrap = $('#task_desc .value');
    var href = location.href.split('/');
    var id = href[5];
    var type = href[3];

    switch (type) {
        case 'tasks': {
            var descWrap = $('#task_desc .value');
            $.get("/get/taskdesc", { id: id }).done(function(data){
                var html = md.render(data);
                descWrap.html(html);
            })
            break
        }
        case 'wiki': {
            var descWrap = $('#wiki_desc .value');
            $.get("/get/wikiarticle", { id: id }).done(function(data){
                var html = md.render(data);
                descWrap.html(html);
            })
            break
        }
    }
});