$(document).ready(function() {
    console.log("ready")
    var md = window.markdownit();
    $('#input_desc').bind('input propertychange', function() {
        var markdown = $('#input_desc').val();
        console.log(markdown)
        var html = md.render(markdown);
        $('#markdown_output').html(html);
    })
})