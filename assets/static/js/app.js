/*
 * MIT License
 *
 * Copyright (c) 2020 Dan <https://github.com/delivrance>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

class Nekobin {
  constructor() {
    this.theme = "dark";

    this.actions = {
      theme: document.getElementById("theme"),
      raw: document.getElementById("raw"),
      save: document.getElementById("save"),
      new: document.getElementById("new")
    };

    CodeMirror.modeURL = "https://cdnjs.cloudflare.com/ajax/libs/codemirror/5.48.4/mode/%N/%N.min.js";

    this.editor = CodeMirror(document.getElementById("content"), {
      placeholder: "Paste code, save and share the link!",
      lineNumbers: true
    });

    this.url = document.getElementById("url");
    this.editor.focus();
  }

  async animateURL() {
    let prevHTML = this.url.innerHTML;

    this.url.classList.add("unclickable");
    this.url.style.opacity = "0";
    await sleep(150);

    this.url.innerHTML = `<i class="fas fa-check"></i> Copied!`;
    this.url.style.opacity = "1";
    await sleep(1500);

    this.url.style.opacity = "0";
    await sleep(150);

    this.url.innerHTML = prevHTML;
    this.url.style.opacity = null;
    this.url.classList.remove("unclickable");
  }

  isContentEmpty() {
    return this.editor.getDoc().getValue().length === 0;
  }

  async switchTheme() {
    let themeEl = this.actions.theme;

    function getProp(prop) {
      return getComputedStyle(document.documentElement).getPropertyValue(prop);
    }

    function setProp(prop, value) {
      document.documentElement.style.setProperty(prop, value);
    }

    if (this.theme === "dark") {
      themeEl.classList.remove("fa-moon");
      themeEl.classList.add("fa-sun");

      setProp("--bg-color", getProp("--bg-light-color"));
      setProp("--bg2-color", getProp("--bg2-light-color"));
      setProp("--main-color", getProp("--main-light-color"));

      setProp("--border-color", getProp("--border-light-color"));
      setProp("--scrollbar-color", getProp("--scrollbar-light-color"));
      setProp("--scrollbar-active-color", getProp("--scrollbar-active-light-color"));

      setProp("--placeholder-color", getProp("--placeholder-light-color"));
      setProp("--linenumber-color", getProp("--linenumber-light-color"));

      this.editor.setOption("theme", "default");

      this.theme = "light"
    } else {
      themeEl.classList.remove("fa-sun");
      themeEl.classList.add("fa-moon");

      setProp("--bg-color", getProp("--bg-dark-color"));
      setProp("--bg2-color", getProp("--bg2-dark-color"));
      setProp("--main-color", getProp("--main-dark-color"));

      setProp("--border-color", getProp("--border-dark-color"));
      setProp("--scrollbar-color", getProp("--scrollbar-dark-color"));
      setProp("--scrollbar-active-color", getProp("--scrollbar-active-dark-color"));

      setProp("--placeholder-color", getProp("--placeholder-dark-color"));
      setProp("--linenumber-color", getProp("--linenumber-dark-color"));

      this.editor.setOption("theme", "darcula");

      this.theme = "dark";
    }

    document.cookie = `theme=${this.theme}`;
  }

  async setup() {
    let key = window.location.pathname;

    this.theme = getCookie("theme") || "dark";
    // Call twice to set the theme got from cookies. Rework.
    await this.switchTheme();
    await this.switchTheme();

    this.url.onclick = async () => {
      copyToClipboard(window.location.href);
      await this.animateURL();
    };

    this.actions.theme.onclick = async () => {
      document.body.style.opacity = "0";
      await sleep(150);
      await this.switchTheme();
      document.body.style.opacity = null;
    };

    this.actions.raw.onclick = () => {
      window.location.href = `/raw${key}`;
    };

    this.actions.new.onclick = () => {
      window.location.href = "/";
    };

    this.actions.save.onclick = async () => {
      this.actions.save.disabled = true;

      let content = this.editor.getDoc().getValue();

      let response = await fetch("/api/documents", {
        method: "POST",
        headers: {"Content-Type": "application/json"},
        body: JSON.stringify({content})
      });

      if (response.ok) {
        let {key} = (await response.json()).result;
        window.location.href = `/${key}`;
      } else {
        let {error} = await response.json();
        this.actions.save.disabled = false;
        alert(`Error: ${error}`);
      }
    };

    this.editor.on("change", () => {
      this.actions.save.disabled = this.isContentEmpty();
    });

    this.editor.setOption("extraKeys", {
      "Ctrl-S": () => this.actions.save.click(),
      "Shift-Ctrl-R": () => this.actions.raw.click(),
      "Ctrl-N": () => this.actions.new.click()
    })
  }

  async load() {
    let path = window.location.pathname;

    if (path === "/") {
      return
    }

    let response = await fetch(`/api/documents${path}`);

    if (response.ok) {
      let {key, content} = (await response.json()).result;

      this.editor.getDoc().setValue(content);
      this.editor.setOption("readOnly", true);

      let mode = CodeMirror.findModeByFileName(path);

      if (mode !== undefined) {
        CodeMirror.autoLoadMode(this.editor, mode.mode);
        this.editor.setOption("mode", mode.mime);
      }

      document.getElementById("content").classList.add("readonly");
      document.title = `nekobin - ${key}`;

      let url = document.getElementById("url");

      url.insertAdjacentText("afterbegin", path);
      url.classList.remove("hidden");

      this.actions.save.disabled = true;

      if (key !== "about") {
        this.actions.raw.disabled = false;
      }
    } else {
      if (response.status === 429) {
        let {error} = await response.json();
        alert(`Error: ${error}`);
      } else {
        window.location.replace("/");
      }
    }
  }
}

// https://www.w3schools.com/js/js_cookies.asp
function getCookie(cname) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');

  for (let i = 0; i < ca.length; i++) {
    let c = ca[i];

    while (c.charAt(0) === ' ') {
      c = c.substring(1);
    }

    if (c.indexOf(name) === 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}

// https://stackoverflow.com/questions/33855641/copy-output-of-a-javascript-variable-to-the-clipboard
const copyToClipboard = text => {
  let dummy = document.createElement("textarea");

  document.body.appendChild(dummy);

  dummy.value = text;
  dummy.select();

  document.execCommand("copy");

  document.body.removeChild(dummy);
};

// https://stackoverflow.com/questions/951021/what-is-the-javascript-version-of-sleep
const sleep = ms => new Promise(r => setTimeout(r, ms));

window.addEventListener("DOMContentLoaded", async () => {
  let nekobin = new Nekobin();

  await nekobin.setup();
  await nekobin.load();
});
