$(document).ready(function() {
    var href = location.href.split('/');
    var type = href[3];
    var id = $('.pagetitle').attr("id");

    getTaskDescOrWikiArticle(type, id);
});