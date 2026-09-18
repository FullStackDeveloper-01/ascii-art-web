// ascii-art-web — small vanilla-JS quality-of-life additions (no libraries, no
// dependencies). Purely front-end polish: none of this touches the Go backend,
// the routes, or the form's required fields/names.

var MAX_LENGTH = 10000;

function updateCharCount() {
    var input = document.getElementById('text');
    var count = document.getElementById('charCount');
    var limitMsg = document.getElementById('limitMsg');
    if (!input || !count) return;

    var len = input.value.length;
    var atMax = len >= MAX_LENGTH;
    count.textContent = len + ' / ' + MAX_LENGTH;

    count.classList.remove('warn', 'danger');
    if (atMax) {
        count.classList.add('danger');
    } else if (len >= MAX_LENGTH * 0.95) {
        count.classList.add('danger');
    } else if (len >= MAX_LENGTH * 0.8) {
        count.classList.add('warn');
    }
    if (limitMsg) limitMsg.hidden = !atMax;
}

function clearText() {
    var input = document.getElementById('text');
    if (!input) return;
    input.value = '';
    input.focus();
    updateCharCount();
}

function copyResult() {
    var pre = document.getElementById('resultText');
    var toast = document.getElementById('copyToast');
    if (!pre) return;

    navigator.clipboard.writeText(pre.innerText).then(function () {
        if (!toast) return;
        toast.classList.add('show');
        window.clearTimeout(copyResult._t);
        copyResult._t = window.setTimeout(function () {
            toast.classList.remove('show');
        }, 1500);
    });
}

document.addEventListener('DOMContentLoaded', function () {
    updateCharCount();
});
