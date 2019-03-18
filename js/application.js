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

    $('#ogp').on('change', function() {
        if(this.checked) {
            $('.ogp-information').show();
            if($('input.ogp[name="title"]').val() === "") {
                fetchOGP();
            }
        } else {
            $('.ogp-information').hide();
        }
    })

    function fetchOGP() {
        var url = $('#url').val();
        if(url === "") {
            return;
        }
        $('input.ogp').val("取得中...");
        $('input.ogp').prop("disabled", true);
        $.ajax({
            type: 'GET',
            dataType: 'json',
            url: "https://ogp.en-courage.com?url=" + url,
            success: function(data){
                $('input.ogp[name="title"]').val(data.title);
                $('input.ogp[name="image"]').val(data.image);
                $('input.ogp[name="description"]').val(data.description);
                $('input.ogp').prop("disabled", false);
            }
        })
    }
})