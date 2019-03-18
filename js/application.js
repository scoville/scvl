$(function(){
    if($('#shortenUrl')[0]) {
        $('#shortenUrl').val(location.origin + "/" + $('#shortenUrl').val());
    }
    $('.copy i').on('click', function(){
        $('#shortenUrl').select();
        document.execCommand("copy");
        $('.copy .hint').css('visibility', 'visible').text('コピー完了');
        setTimeout(function(){
            $('.copy .hint').css('visibility', '');
        }, 3000)
    });
})