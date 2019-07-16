$(document).ready(function() {
    var md = window.markdownit();
    $('#input_desc').bind('input propertychange', function() {
        var markdown = $('#input_desc').val();
        var html = md.render(markdown);
        $('#markdown_output').html(html);
    })
})