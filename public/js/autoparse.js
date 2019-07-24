$(document).ready(function() {
    $('#input_desc').bind('input propertychange', function() {
        var markdown = $('#input_desc').val();
        MDParse(markdown, $('#markdown_output'))
    })
})