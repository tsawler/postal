package postal

var bootstrapPartial = `
{{define "bootstrap-css"}}
    /*!
    * Bootstrap v4.5.0 (https://getbootstrap.com/)
    * Copyright 2011-2020 The Bootstrap Authors
    * Copyright 2011-2020 Twitter, Inc.
    * Licensed under MIT (https://github.com/twbs/bootstrap/blob/master/LICENSE)
    */
    :root {
    --blue: #007bff;
    --indigo: #6610f2;
    --purple: #6f42c1;
    --pink: #e83e8c;
    --red: #dc3545;
    --orange: #fd7e14;
    --yellow: #ffc107;
    --green: #28a745;
    --teal: #20c997;
    --cyan: #17a2b8;
    --white: #fff;
    --gray: #6c757d;
    --gray-dark: #343a40;
    --primary: #007bff;
    --secondary: #6c757d;
    --success: #28a745;
    --info: #17a2b8;
    --warning: #ffc107;
    --danger: #dc3545;
    --light: #f8f9fa;
    --dark: #343a40;
    --breakpoint-xs: 0;
    --breakpoint-sm: 576px;
    --breakpoint-md: 768px;
    --breakpoint-lg: 992px;
    --breakpoint-xl: 1200px;
    --font-family-sans-serif: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
    --font-family-monospace: SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    }

    *,
    *::before,
    *::after {
    box-sizing: border-box;
    }

    html {
    font-family: sans-serif;
    line-height: 1.15;
    -webkit-text-size-adjust: 100%;
    -webkit-tap-highlight-color: rgba(0, 0, 0, 0);
    }

    figcaption, figure, footer {
    display: block;
    }

    body {
    margin: 0;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif, "Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
    font-size: 1rem;
    font-weight: 400;
    line-height: 1.5;
    color: #212529;
    text-align: left;
    background-color: #fff;
    }

    [tabindex="-1"]:focus:not(:focus-visible) {
    outline: 0 !important;
    }

    h1, h2, h3, h4, h5 {
    margin-top: 0;
    margin-bottom: 0.5rem;
    }

    p {
    margin-top: 0;
    margin-bottom: 1rem;
    }

    abbr[title] {
    text-decoration: underline;
    -webkit-text-decoration: underline dotted;
    text-decoration: underline dotted;
    cursor: help;
    border-bottom: 0;
    -webkit-text-decoration-skip-ink: none;
    text-decoration-skip-ink: none;
    }


    ul,
    dl {
    margin-top: 0;
    margin-bottom: 1rem;
    }


    ul ul {
    margin-bottom: 0;
    }

    dt {
    font-weight: 700;
    }

    dd {
    margin-bottom: .5rem;
    margin-left: 0;
    }

    blockquote {
    margin: 0 0 1rem;
    }


    strong {
    font-weight: bolder;
    }

    small {
    font-size: 80%;
    }

    a {
    color: #007bff;
    text-decoration: none;
    background-color: transparent;
    }

    a:hover {
    color: #0056b3;
    text-decoration: underline;
    }

    a:not([href]) {
    color: inherit;
    text-decoration: none;
    }

    a:not([href]):hover {
    color: inherit;
    text-decoration: none;
    }


    code {
    font-family: SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace;
    font-size: 1em;
    }

    figure {
    margin: 0 0 1rem;
    }

    img {
    vertical-align: middle;
    border-style: none;
    }

    table {
    border-collapse: collapse;
    }

    caption {
    padding-top: 0.75rem;
    padding-bottom: 0.75rem;
    color: #6c757d;
    text-align: left;
    caption-side: bottom;
    }

    th {
    text-align: inherit;
    }

    h1, h2, h3, h4, h5,
    .h1, .h2, .h3, .h4, .h5 {
    margin-bottom: 0.5rem;
    font-weight: 500;
    line-height: 1.2;
    }

    h1, .h1 {
    font-size: 2.5rem;
    }

    h2, .h2 {
    font-size: 2rem;
    }

    h3, .h3 {
    font-size: 1.75rem;
    }

    h4, .h4 {
    font-size: 1.5rem;
    }

    h5, .h5 {
    font-size: 1.25rem;
    }

    .lead {
    font-size: 1.25rem;
    font-weight: 300;
    }

    .display-1 {
    font-size: 6rem;
    font-weight: 300;
    line-height: 1.2;
    }

    .display-2 {
    font-size: 5.5rem;
    font-weight: 300;
    line-height: 1.2;
    }

    .display-3 {
    font-size: 4.5rem;
    font-weight: 300;
    line-height: 1.2;
    }

    .display-4 {
    font-size: 3.5rem;
    font-weight: 300;
    line-height: 1.2;
    }

    small,
    .small {
    font-size: 80%;
    font-weight: 400;
    }

    mark,
    .mark {
    padding: 0.2em;
    background-color: #fcf8e3;
    }

    .list-unstyled {
    padding-left: 0;
    list-style: none;
    }

    .list-inline {
    padding-left: 0;
    list-style: none;
    }

    .list-inline-item {
    display: inline-block;
    }

    .list-inline-item:not(:last-child) {
    margin-right: 0.5rem;
    }

    .initialism {
    font-size: 90%;
    text-transform: uppercase;
    }

    .blockquote {
    margin-bottom: 1rem;
    font-size: 1.25rem;
    }

    .blockquote-footer {
    display: block;
    font-size: 80%;
    color: #6c757d;
    }

    .blockquote-footer::before {
    content: "\2014\00A0";
    }

    .img-fluid {
    max-width: 100%;
    height: auto;
    }

    .img-thumbnail {
    padding: 0.25rem;
    background-color: #fff;
    border: 1px solid #dee2e6;
    border-radius: 0.25rem;
    max-width: 100%;
    height: auto;
    }

    .figure {
    display: inline-block;
    }

    .figure-img {
    margin-bottom: 0.5rem;
    line-height: 1;
    }

    .figure-caption {
    font-size: 90%;
    color: #6c757d;
    }

    code {
    font-size: 87.5%;
    color: #e83e8c;
    word-wrap: break-word;
    }

    a > code {
    color: inherit;
    }

    .container {
    width: 100%;
    padding-right: 15px;
    padding-left: 15px;
    margin-right: auto;
    margin-left: auto;
    }

    @media (min-width: 576px) {
    .container {
    max-width: 540px;
    }
    }

    @media (min-width: 768px) {
    .container {
    max-width: 720px;
    }
    }

    @media (min-width: 992px) {
    .container {
    max-width: 960px;
    }
    }

    @media (min-width: 1200px) {
    .container {
    max-width: 1140px;
    }
    }

    @media (min-width: 576px) {
    .container {
    max-width: 540px;
    }
    }

    @media (min-width: 768px) {
    .container {
    max-width: 720px;
    }
    }

    @media (min-width: 992px) {
    .container {
    max-width: 960px;
    }
    }

    @media (min-width: 1200px) {
    .container {
    max-width: 1140px;
    }
    }

    .row {
    display: -ms-flexbox;
    display: flex;
    -ms-flex-wrap: wrap;
    flex-wrap: wrap;
    margin-right: -15px;
    margin-left: -15px;
    }

    .col-1, .col-2, .col-4, .col-6, .col-8, .col-10, .col-12, .col, .col-sm-3, .col-sm-4, .col-sm-8, .col-sm-9, .col-md-1, .col-md-2, .col-md-4, .col-md-6, .col-md-8, .col-md-10, .col-md-12, .col-lg-1, .col-lg-2, .col-lg-4, .col-lg-6, .col-lg-8, .col-lg-10, .col-lg-12 {
    position: relative;
    width: 100%;
    padding-right: 15px;
    padding-left: 15px;
    }

    .col {
    -ms-flex-preferred-size: 0;
    flex-basis: 0;
    -ms-flex-positive: 1;
    flex-grow: 1;
    min-width: 0;
    max-width: 100%;
    }

    .col-1 {
    -ms-flex: 0 0 8.333333%;
    flex: 0 0 8.333333%;
    max-width: 8.333333%;
    }

    .col-2 {
    -ms-flex: 0 0 16.666667%;
    flex: 0 0 16.666667%;
    max-width: 16.666667%;
    }

    .col-4 {
    -ms-flex: 0 0 33.333333%;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;
    }

    .col-6 {
    -ms-flex: 0 0 50%;
    flex: 0 0 50%;
    max-width: 50%;
    }

    .col-8 {
    -ms-flex: 0 0 66.666667%;
    flex: 0 0 66.666667%;
    max-width: 66.666667%;
    }

    .col-10 {
    -ms-flex: 0 0 83.333333%;
    flex: 0 0 83.333333%;
    max-width: 83.333333%;
    }

    .col-12 {
    -ms-flex: 0 0 100%;
    flex: 0 0 100%;
    max-width: 100%;
    }

    @media (min-width: 576px) {
    .col-sm-3 {
    -ms-flex: 0 0 25%;
    flex: 0 0 25%;
    max-width: 25%;
    }
    .col-sm-4 {
    -ms-flex: 0 0 33.333333%;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;
    }
    .col-sm-8 {
    -ms-flex: 0 0 66.666667%;
    flex: 0 0 66.666667%;
    max-width: 66.666667%;
    }
    .col-sm-9 {
    -ms-flex: 0 0 75%;
    flex: 0 0 75%;
    max-width: 75%;
    }
    }

    @media (min-width: 768px) {
    .col-md-1 {
    -ms-flex: 0 0 8.333333%;
    flex: 0 0 8.333333%;
    max-width: 8.333333%;
    }
    .col-md-2 {
    -ms-flex: 0 0 16.666667%;
    flex: 0 0 16.666667%;
    max-width: 16.666667%;
    }
    .col-md-4 {
    -ms-flex: 0 0 33.333333%;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;
    }
    .col-md-6 {
    -ms-flex: 0 0 50%;
    flex: 0 0 50%;
    max-width: 50%;
    }
    .col-md-8 {
    -ms-flex: 0 0 66.666667%;
    flex: 0 0 66.666667%;
    max-width: 66.666667%;
    }
    .col-md-10 {
    -ms-flex: 0 0 83.333333%;
    flex: 0 0 83.333333%;
    max-width: 83.333333%;
    }
    .col-md-12 {
    -ms-flex: 0 0 100%;
    flex: 0 0 100%;
    max-width: 100%;
    }
    }

    @media (min-width: 992px) {
    .col-lg-1 {
    -ms-flex: 0 0 8.333333%;
    flex: 0 0 8.333333%;
    max-width: 8.333333%;
    }
    .col-lg-2 {
    -ms-flex: 0 0 16.666667%;
    flex: 0 0 16.666667%;
    max-width: 16.666667%;
    }
    .col-lg-4 {
    -ms-flex: 0 0 33.333333%;
    flex: 0 0 33.333333%;
    max-width: 33.333333%;
    }
    .col-lg-6 {
    -ms-flex: 0 0 50%;
    flex: 0 0 50%;
    max-width: 50%;
    }
    .col-lg-8 {
    -ms-flex: 0 0 66.666667%;
    flex: 0 0 66.666667%;
    max-width: 66.666667%;
    }
    .col-lg-10 {
    -ms-flex: 0 0 83.333333%;
    flex: 0 0 83.333333%;
    max-width: 83.333333%;
    }
    .col-lg-12 {
    -ms-flex: 0 0 100%;
    flex: 0 0 100%;
    max-width: 100%;
    }
    }

    .table {
    width: 100%;
    margin-bottom: 1rem;
    color: #212529;
    }

    .table th,
    .table td {
    padding: 0.75rem;
    vertical-align: top;
    border-top: 1px solid #dee2e6;
    }

    .table thead th {
    vertical-align: bottom;
    border-bottom: 2px solid #dee2e6;
    }

    .table tbody + tbody {
    border-top: 2px solid #dee2e6;
    }

    .table-striped tbody tr:nth-of-type(odd) {
    background-color: rgba(0, 0, 0, 0.05);
    }

    .btn {
    display: inline-block;
    font-weight: 400;
    color: #212529;
    text-align: center;
    vertical-align: middle;
    -webkit-user-select: none;
    -moz-user-select: none;
    -ms-user-select: none;
    user-select: none;
    background-color: transparent;
    border: 1px solid transparent;
    padding: 0.375rem 0.75rem;
    font-size: 1rem;
    line-height: 1.5;
    border-radius: 0.25rem;
    transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
    }

    @media (prefers-reduced-motion: reduce) {
    .btn {
    transition: none;
    }
    }

    .btn:hover {
    color: #212529;
    text-decoration: none;
    }

    .btn:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
    }

    .btn.disabled, .btn:disabled {
    opacity: 0.65;
    }

    .btn:not(:disabled):not(.disabled) {
    cursor: pointer;
    }

    a.btn.disabled {
    pointer-events: none;
    }

    .btn-primary {
    color: #fff;
    background-color: #007bff;
    border-color: #007bff;
    }

    .btn-primary:hover {
    color: #fff;
    background-color: #0069d9;
    border-color: #0062cc;
    }

    .btn-primary:focus {
    color: #fff;
    background-color: #0069d9;
    border-color: #0062cc;
    box-shadow: 0 0 0 0.2rem rgba(38, 143, 255, 0.5);
    }

    .btn-primary.disabled, .btn-primary:disabled {
    color: #fff;
    background-color: #007bff;
    border-color: #007bff;
    }

    .btn-primary:not(:disabled):not(.disabled):active, .btn-primary:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #0062cc;
    border-color: #005cbf;
    }

    .btn-primary:not(:disabled):not(.disabled):active:focus, .btn-primary:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(38, 143, 255, 0.5);
    }

    .btn-secondary {
    color: #fff;
    background-color: #6c757d;
    border-color: #6c757d;
    }

    .btn-secondary:hover {
    color: #fff;
    background-color: #5a6268;
    border-color: #545b62;
    }

    .btn-secondary:focus {
    color: #fff;
    background-color: #5a6268;
    border-color: #545b62;
    box-shadow: 0 0 0 0.2rem rgba(130, 138, 145, 0.5);
    }

    .btn-secondary.disabled, .btn-secondary:disabled {
    color: #fff;
    background-color: #6c757d;
    border-color: #6c757d;
    }

    .btn-secondary:not(:disabled):not(.disabled):active, .btn-secondary:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #545b62;
    border-color: #4e555b;
    }

    .btn-secondary:not(:disabled):not(.disabled):active:focus, .btn-secondary:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(130, 138, 145, 0.5);
    }

    .btn-success {
    color: #fff;
    background-color: #28a745;
    border-color: #28a745;
    }

    .btn-success:hover {
    color: #fff;
    background-color: #218838;
    border-color: #1e7e34;
    }

    .btn-success:focus {
    color: #fff;
    background-color: #218838;
    border-color: #1e7e34;
    box-shadow: 0 0 0 0.2rem rgba(72, 180, 97, 0.5);
    }

    .btn-success.disabled, .btn-success:disabled {
    color: #fff;
    background-color: #28a745;
    border-color: #28a745;
    }

    .btn-success:not(:disabled):not(.disabled):active, .btn-success:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #1e7e34;
    border-color: #1c7430;
    }

    .btn-success:not(:disabled):not(.disabled):active:focus, .btn-success:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(72, 180, 97, 0.5);
    }

    .btn-info {
    color: #fff;
    background-color: #17a2b8;
    border-color: #17a2b8;
    }

    .btn-info:hover {
    color: #fff;
    background-color: #138496;
    border-color: #117a8b;
    }

    .btn-info:focus {
    color: #fff;
    background-color: #138496;
    border-color: #117a8b;
    box-shadow: 0 0 0 0.2rem rgba(58, 176, 195, 0.5);
    }

    .btn-info.disabled, .btn-info:disabled {
    color: #fff;
    background-color: #17a2b8;
    border-color: #17a2b8;
    }

    .btn-info:not(:disabled):not(.disabled):active, .btn-info:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #117a8b;
    border-color: #10707f;
    }

    .btn-info:not(:disabled):not(.disabled):active:focus, .btn-info:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(58, 176, 195, 0.5);
    }

    .btn-warning {
    color: #212529;
    background-color: #ffc107;
    border-color: #ffc107;
    }

    .btn-warning:hover {
    color: #212529;
    background-color: #e0a800;
    border-color: #d39e00;
    }

    .btn-warning:focus {
    color: #212529;
    background-color: #e0a800;
    border-color: #d39e00;
    box-shadow: 0 0 0 0.2rem rgba(222, 170, 12, 0.5);
    }

    .btn-warning.disabled, .btn-warning:disabled {
    color: #212529;
    background-color: #ffc107;
    border-color: #ffc107;
    }

    .btn-warning:not(:disabled):not(.disabled):active, .btn-warning:not(:disabled):not(.disabled).active {
    color: #212529;
    background-color: #d39e00;
    border-color: #c69500;
    }

    .btn-warning:not(:disabled):not(.disabled):active:focus, .btn-warning:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(222, 170, 12, 0.5);
    }

    .btn-danger {
    color: #fff;
    background-color: #dc3545;
    border-color: #dc3545;
    }

    .btn-danger:hover {
    color: #fff;
    background-color: #c82333;
    border-color: #bd2130;
    }

    .btn-danger:focus {
    color: #fff;
    background-color: #c82333;
    border-color: #bd2130;
    box-shadow: 0 0 0 0.2rem rgba(225, 83, 97, 0.5);
    }

    .btn-danger.disabled, .btn-danger:disabled {
    color: #fff;
    background-color: #dc3545;
    border-color: #dc3545;
    }

    .btn-danger:not(:disabled):not(.disabled):active, .btn-danger:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #bd2130;
    border-color: #b21f2d;
    }

    .btn-danger:not(:disabled):not(.disabled):active:focus, .btn-danger:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(225, 83, 97, 0.5);
    }

    .btn-light {
    color: #212529;
    background-color: #f8f9fa;
    border-color: #f8f9fa;
    }

    .btn-light:hover {
    color: #212529;
    background-color: #e2e6ea;
    border-color: #dae0e5;
    }

    .btn-light:focus {
    color: #212529;
    background-color: #e2e6ea;
    border-color: #dae0e5;
    box-shadow: 0 0 0 0.2rem rgba(216, 217, 219, 0.5);
    }

    .btn-light.disabled, .btn-light:disabled {
    color: #212529;
    background-color: #f8f9fa;
    border-color: #f8f9fa;
    }

    .btn-light:not(:disabled):not(.disabled):active, .btn-light:not(:disabled):not(.disabled).active {
    color: #212529;
    background-color: #dae0e5;
    border-color: #d3d9df;
    }

    .btn-light:not(:disabled):not(.disabled):active:focus, .btn-light:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(216, 217, 219, 0.5);
    }

    .btn-dark {
    color: #fff;
    background-color: #343a40;
    border-color: #343a40;
    }

    .btn-dark:hover {
    color: #fff;
    background-color: #23272b;
    border-color: #1d2124;
    }

    .btn-dark:focus {
    color: #fff;
    background-color: #23272b;
    border-color: #1d2124;
    box-shadow: 0 0 0 0.2rem rgba(82, 88, 93, 0.5);
    }

    .btn-dark.disabled, .btn-dark:disabled {
    color: #fff;
    background-color: #343a40;
    border-color: #343a40;
    }

    .btn-dark:not(:disabled):not(.disabled):active, .btn-dark:not(:disabled):not(.disabled).active {
    color: #fff;
    background-color: #1d2124;
    border-color: #171a1d;
    }

    .btn-dark:not(:disabled):not(.disabled):active:focus, .btn-dark:not(:disabled):not(.disabled).active:focus {
    box-shadow: 0 0 0 0.2rem rgba(82, 88, 93, 0.5);
    }

    .card {
    position: relative;
    display: -ms-flexbox;
    display: flex;
    -ms-flex-direction: column;
    flex-direction: column;
    min-width: 0;
    word-wrap: break-word;
    background-color: #fff;
    background-clip: border-box;
    border: 1px solid rgba(0, 0, 0, 0.125);
    border-radius: 0.25rem;
    }

    .card > .list-group {
    border-top: inherit;
    border-bottom: inherit;
    }

    .card > .list-group:first-child {
    border-top-width: 0;
    border-top-left-radius: calc(0.25rem - 1px);
    border-top-right-radius: calc(0.25rem - 1px);
    }

    .card > .list-group:last-child {
    border-bottom-width: 0;
    border-bottom-right-radius: calc(0.25rem - 1px);
    border-bottom-left-radius: calc(0.25rem - 1px);
    }

    .card-body {
    -ms-flex: 1 1 auto;
    flex: 1 1 auto;
    min-height: 1px;
    padding: 1.25rem;
    }

    .card-title {
    margin-bottom: 0.75rem;
    }

    .card-text:last-child {
    margin-bottom: 0;
    }


    .card-img-top {
    -ms-flex-negative: 0;
    flex-shrink: 0;
    width: 100%;
    }


    .card-img-top {
    border-top-left-radius: calc(0.25rem - 1px);
    border-top-right-radius: calc(0.25rem - 1px);
    }

    .badge {
    display: inline-block;
    padding: 0.25em 0.4em;
    font-size: 75%;
    font-weight: 700;
    line-height: 1;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: 0.25rem;
    transition: color 0.15s ease-in-out, background-color 0.15s ease-in-out, border-color 0.15s ease-in-out, box-shadow 0.15s ease-in-out;
    }

    @media (prefers-reduced-motion: reduce) {
    .badge {
    transition: none;
    }
    }

    a.badge:hover, a.badge:focus {
    text-decoration: none;
    }

    .badge:empty {
    display: none;
    }

    .btn .badge {
    position: relative;
    top: -1px;
    }

    .badge-primary {
    color: #fff;
    background-color: #007bff;
    }

    a.badge-primary:hover, a.badge-primary:focus {
    color: #fff;
    background-color: #0062cc;
    }

    a.badge-primary:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.5);
    }

    .badge-secondary {
    color: #fff;
    background-color: #6c757d;
    }

    a.badge-secondary:hover, a.badge-secondary:focus {
    color: #fff;
    background-color: #545b62;
    }

    a.badge-secondary:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(108, 117, 125, 0.5);
    }

    .badge-success {
    color: #fff;
    background-color: #28a745;
    }

    a.badge-success:hover, a.badge-success:focus {
    color: #fff;
    background-color: #1e7e34;
    }

    a.badge-success:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(40, 167, 69, 0.5);
    }

    .badge-info {
    color: #fff;
    background-color: #17a2b8;
    }

    a.badge-info:hover, a.badge-info:focus {
    color: #fff;
    background-color: #117a8b;
    }

    a.badge-info:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(23, 162, 184, 0.5);
    }

    .badge-warning {
    color: #212529;
    background-color: #ffc107;
    }

    a.badge-warning:hover, a.badge-warning:focus {
    color: #212529;
    background-color: #d39e00;
    }

    a.badge-warning:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(255, 193, 7, 0.5);
    }

    .badge-danger {
    color: #fff;
    background-color: #dc3545;
    }

    a.badge-danger:hover, a.badge-danger:focus {
    color: #fff;
    background-color: #bd2130;
    }

    a.badge-danger:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(220, 53, 69, 0.5);
    }

    .badge-light {
    color: #212529;
    background-color: #f8f9fa;
    }

    a.badge-light:hover, a.badge-light:focus {
    color: #212529;
    background-color: #dae0e5;
    }

    a.badge-light:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(248, 249, 250, 0.5);
    }

    .badge-dark {
    color: #fff;
    background-color: #343a40;
    }

    a.badge-dark:hover, a.badge-dark:focus {
    color: #fff;
    background-color: #1d2124;
    }

    a.badge-dark:focus {
    outline: 0;
    box-shadow: 0 0 0 0.2rem rgba(52, 58, 64, 0.5);
    }

    .alert {
    position: relative;
    padding: 0.75rem 1.25rem;
    margin-bottom: 1rem;
    border: 1px solid transparent;
    border-radius: 0.25rem;
    }

    .alert-link {
    font-weight: 700;
    }

    .alert-primary {
    color: #004085;
    background-color: #cce5ff;
    border-color: #b8daff;
    }

    .alert-primary .alert-link {
    color: #002752;
    }

    .alert-secondary {
    color: #383d41;
    background-color: #e2e3e5;
    border-color: #d6d8db;
    }

    .alert-secondary .alert-link {
    color: #202326;
    }

    .alert-success {
    color: #155724;
    background-color: #d4edda;
    border-color: #c3e6cb;
    }

    .alert-success .alert-link {
    color: #0b2e13;
    }

    .alert-info {
    color: #0c5460;
    background-color: #d1ecf1;
    border-color: #bee5eb;
    }

    .alert-info .alert-link {
    color: #062c33;
    }

    .alert-warning {
    color: #856404;
    background-color: #fff3cd;
    border-color: #ffeeba;
    }

    .alert-warning .alert-link {
    color: #533f03;
    }

    .alert-danger {
    color: #721c24;
    background-color: #f8d7da;
    border-color: #f5c6cb;
    }

    .alert-danger .alert-link {
    color: #491217;
    }

    .alert-light {
    color: #818182;
    background-color: #fefefe;
    border-color: #fdfdfe;
    }

    .alert-light .alert-link {
    color: #686868;
    }

    .alert-dark {
    color: #1b1e21;
    background-color: #d6d8d9;
    border-color: #c6c8ca;
    }

    .alert-dark .alert-link {
    color: #040505;
    }

    .list-group {
    display: -ms-flexbox;
    display: flex;
    -ms-flex-direction: column;
    flex-direction: column;
    padding-left: 0;
    margin-bottom: 0;
    border-radius: 0.25rem;
    }

    .list-group-item-action {
    width: 100%;
    color: #495057;
    text-align: inherit;
    }

    .list-group-item-action:hover, .list-group-item-action:focus {
    z-index: 1;
    color: #495057;
    text-decoration: none;
    background-color: #f8f9fa;
    }

    .list-group-item-action:active {
    color: #212529;
    background-color: #e9ecef;
    }

    .list-group-item {
    position: relative;
    display: block;
    padding: 0.75rem 1.25rem;
    background-color: #fff;
    border: 1px solid rgba(0, 0, 0, 0.125);
    }

    .list-group-item:first-child {
    border-top-left-radius: inherit;
    border-top-right-radius: inherit;
    }

    .list-group-item:last-child {
    border-bottom-right-radius: inherit;
    border-bottom-left-radius: inherit;
    }

    .list-group-item.disabled, .list-group-item:disabled {
    color: #6c757d;
    pointer-events: none;
    background-color: #fff;
    }

    .list-group-item.active {
    z-index: 2;
    color: #fff;
    background-color: #007bff;
    border-color: #007bff;
    }

    .list-group-item + .list-group-item {
    border-top-width: 0;
    }

    .list-group-item + .list-group-item.active {
    margin-top: -1px;
    border-top-width: 1px;
    }

    .bg-dark {
    background-color: #343a40 !important;
    }

    a.bg-dark:hover, a.bg-dark:focus {
    background-color: #1d2124 !important;
    }

    .rounded {
    border-radius: 0.25rem !important;
    }

    .d-block {
    display: block !important;
    }

    .float-left {
    float: left !important;
    }

    .float-right {
    float: right !important;
    }

    .mb-0 {
    margin-bottom: 0 !important;
    }

    .mt-1 {
    margin-top: 0.25rem !important;
    }

    .mr-1 {
    margin-right: 0.25rem !important;
    }

    .mb-1 {
    margin-bottom: 0.25rem !important;
    }

    .ml-1 {
    margin-left: 0.25rem !important;
    }

    .mt-2 {
    margin-top: 0.5rem !important;
    }

    .mr-2 {
    margin-right: 0.5rem !important;
    }

    .mb-2 {
    margin-bottom: 0.5rem !important;
    }

    .ml-2 {
    margin-left: 0.5rem !important;
    }

    .mt-3 {
    margin-top: 1rem !important;
    }

    .mr-3 {
    margin-right: 1rem !important;
    }

    .mb-3 {
    margin-bottom: 1rem !important;
    }

    .ml-3 {
    margin-left: 1rem !important;
    }

    .mt-4 {
    margin-top: 1.5rem !important;
    }

    .mr-4 {
    margin-right: 1.5rem !important;
    }

    .mb-4 {
    margin-bottom: 1.5rem !important;
    }

    .ml-4 {
    margin-left: 1.5rem !important;
    }

    .mt-5 {
    margin-top: 3rem !important;
    }

    .mr-5 {
    margin-right: 3rem !important;
    }

    .mb-5 {
    margin-bottom: 3rem !important;
    }

    .ml-5 {
    margin-left: 3rem !important;
    }

    .pt-1 {
    padding-top: 0.25rem !important;
    }

    .pr-1 {
    padding-right: 0.25rem !important;
    }

    .pb-1 {
    padding-bottom: 0.25rem !important;
    }

    .pl-1 {
    padding-left: 0.25rem !important;
    }

    .pt-2 {
    padding-top: 0.5rem !important;
    }

    .pr-2 {
    padding-right: 0.5rem !important;
    }

    .pb-2 {
    padding-bottom: 0.5rem !important;
    }

    .pl-2 {
    padding-left: 0.5rem !important;
    }

    .pt-3 {
    padding-top: 1rem !important;
    }

    .pr-3 {
    padding-right: 1rem !important;
    }

    .pb-3 {
    padding-bottom: 1rem !important;
    }

    .pl-3 {
    padding-left: 1rem !important;
    }

    .pt-4 {
    padding-top: 1.5rem !important;
    }

    .pr-4 {
    padding-right: 1.5rem !important;
    }

    .pb-4 {
    padding-bottom: 1.5rem !important;
    }

    .pl-4 {
    padding-left: 1.5rem !important;
    }

    .pt-5 {
    padding-top: 3rem !important;
    }

    .pr-5 {
    padding-right: 3rem !important;
    }

    .pb-5 {
    padding-bottom: 3rem !important;
    }

    .pl-5 {
    padding-left: 3rem !important;
    }


    .mx-auto {
    margin-right: auto !important;
    }


    .mx-auto {
    margin-left: auto !important;
    }

    .text-wrap {
    white-space: normal !important;
    }

    .text-truncate {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    }

    .text-left {
    text-align: left !important;
    }

    .text-right {
    text-align: right !important;
    }

    .text-center {
    text-align: center !important;
    }

    @media (min-width: 576px) {
    .text-sm-left {
    text-align: left !important;
    }
    }

    @media (min-width: 768px) {
    .text-md-left {
    text-align: left !important;
    }
    }

    @media (min-width: 992px) {
    .text-lg-left {
    text-align: left !important;
    }
    }

    @media (min-width: 1200px) {
    .text-xl-left {
    text-align: left !important;
    }
    }

    .font-weight-light {
    font-weight: 300 !important;
    }

    .font-weight-lighter {
    font-weight: lighter !important;
    }

    .font-weight-normal {
    font-weight: 400 !important;
    }

    .font-weight-bold {
    font-weight: 700 !important;
    }

    .font-weight-bolder {
    font-weight: bolder !important;
    }

    .font-italic {
    font-style: italic !important;
    }

    .text-white {
    color: #fff !important;
    }

    .text-primary {
    color: #007bff !important;
    }

    a.text-primary:hover, a.text-primary:focus {
    color: #0056b3 !important;
    }

    .text-secondary {
    color: #6c757d !important;
    }

    a.text-secondary:hover, a.text-secondary:focus {
    color: #494f54 !important;
    }

    .text-success {
    color: #28a745 !important;
    }

    a.text-success:hover, a.text-success:focus {
    color: #19692c !important;
    }

    .text-info {
    color: #17a2b8 !important;
    }

    a.text-info:hover, a.text-info:focus {
    color: #0f6674 !important;
    }

    .text-warning {
    color: #ffc107 !important;
    }

    a.text-warning:hover, a.text-warning:focus {
    color: #ba8b00 !important;
    }

    .text-danger {
    color: #dc3545 !important;
    }

    a.text-danger:hover, a.text-danger:focus {
    color: #a71d2a !important;
    }

    .text-light {
    color: #f8f9fa !important;
    }

    a.text-light:hover, a.text-light:focus {
    color: #cbd3da !important;
    }

    .text-dark {
    color: #343a40 !important;
    }

    a.text-dark:hover, a.text-dark:focus {
    color: #121416 !important;
    }

    .text-body {
    color: #212529 !important;
    }

    .text-muted {
    color: #6c757d !important;
    }

    .text-black-50 {
    color: rgba(0, 0, 0, 0.5) !important;
    }

    .text-white-50 {
    color: rgba(255, 255, 255, 0.5) !important;
    }

    @media print {
    *,
    *::before,
    *::after {
    text-shadow: none !important;
    box-shadow: none !important;
    }
    a:not(.btn) {
    text-decoration: underline;
    }
    abbr[title]::after {
    content: " (" attr(title) ")";
    }

    blockquote {
    border: 1px solid #adb5bd;
    page-break-inside: avoid;
    }
    thead {
    display: table-header-group;
    }
    tr,
    img {
    page-break-inside: avoid;
    }
    p,
    h2,
    h3 {
    orphans: 3;
    widows: 3;
    }
    h2,
    h3 {
    page-break-after: avoid;
    }
    @page {
    size: a3;
    }
    body {
    min-width: 992px !important;
    }
    .container {
    min-width: 992px !important;
    }
    .badge {
    border: 1px solid #000;
    }
    .table {
    border-collapse: collapse !important;
    }
    .table td,
    .table th {
    background-color: #fff !important;
    }
    }
    /*# sourceMappingURL=bootstrap.css.map */
{{end}}
`
