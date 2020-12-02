$(function(){
    if($('#shortenUrl')[0]) {
        $('#shortenUrl').val(location.origin + "/" + $('#shortenUrl').val());
    }
    $('.copy i').on('click', function(){
        $('.copy-target').select();
        document.execCommand("copy");
        $('.copy .hint').css('visibility', 'visible').text('コピー完了');
        setTimeout(function(){
            $('.copy .hint').css('visibility', '');
        }, 3000)
    });

    // UTM
    if($('#url')[0]) {
        var params = getParams($('#url').val());
        if(params["utm_source"]) {
            $('input[name="utm_source"]').val(params["utm_source"]);
        }
        if(params["utm_medium"]) {
            $('select[name="utm_medium"]').val(params["utm_medium"]);
        }
        if(params["utm_campaign"]) {
            $('input[name="utm_campaign"]').val(params["utm_campaign"]);
        }
    }
    if($('#utm')[0] && $('#utm')[0].checked) {
        $('.utm-information').show();
    }
    $('#utm').on('change', function() {
        if(this.checked) {
            $('.utm-information').show();
        } else {
            $('.utm-information').hide();
        }
    });
    $('input[name="url"],input[name="utm_source"],select[name="utm_medium"],input[name="utm_campaign"]').on('change', function() {
        buildURL();
    });

    // OGP
    if($('#ogp')[0] && $('#ogp')[0].checked) {
        $('.ogp-information').show();
    }
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

    // Email
    if($('#email')[0] && $('#email')[0].checked) {
        $('.email-information').show();
    }
    $('#email').on('change', function() {
        if(this.checked) {
            $('.email-information').show();
        } else {
            $('.email-information').hide();
        }
    })

    function getParams(url) {
        var params = {};
        if(url.split("?").length <= 1) {
            return params;
        }
        var q = url.split("?")[1];
        for(var i = 0; i < q.split("&").length; i++) {
            var param = q.split("&")[i].split("=");
            if(param.length > 1){
                params[param[0]] = param[1];
            }
        }
        return params;
    }

    function buildURL() {
        var url = $('#url').val();
        if(url === "") {
            return;
        }
        var params = getParams(url);
        var source = $('input[name="utm_source"]').val()
        if(source) {
            params["utm_source"] = source;
        }
        var medium = $('select[name="utm_medium"]').val();
        if(medium === "-") {
            delete params["utm_medium"];
        } else if(medium) {
            params["utm_medium"] = medium;
        }
        var campaign = $('input[name="utm_campaign"]').val();
        if(campaign) {
            params["utm_campaign"] = campaign;
        }
        url = url.split("?")[0];
        var values = [];
        for(var i = 0; i < Object.keys(params).length; i++) {
            var key = Object.keys(params)[i];
            var val = params[key];
            if(val !== "") {
                values.push(key + "=" + val);
            }
        }
        if(values.length > 0) {
            url = url + "?" + values.join("&");
        }
        $('#url').val(url);
    }

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

    $('.api-key-show').click(function(){
        $('.api-key-show').hide();
        $('.api-key-container').toggle();
    })
    $('#publish-api-key-form').submit(function(){
        return confirm('再発行するとこれまで使用していたAPI Keyは使用できなくなります。\nAPI Keyを再発行しますか？');
    })
})