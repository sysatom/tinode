# Localizing Tinode

**IMPORTANT!** Please use `devel` branches for translations.

## Server

The server sends emails to users upon creation of a new account and when the user requests to reset the password:

* [/server/templ/email-validation-en.templ](../server/templ/email-validation-en.templ)
* [/server/templ/email-password-reset-en.templ](../server/templ/email-password-reset-en.templ)

Create a copy of the files naming them `email-password-reset-XX.teml` and `email-validation-XX.templ` where `XX` is the [ISO-631-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) code of the new language. Translate the content and send a pull request with the new files. If you don't know how to create a pull request then just sent the translated files in any way you can.


## Webapp

The translations are located in two places: [/src/i18n/](https://github.com/tinode/webapp/tree/devel/src/i18n/) and [/service-worker.js](https://github.com/tinode/webapp/blob/devel/service-worker.js#L11).

In order to add a translation, copy `/src/i18n/en.json` to a file named `/src/i18n/XX.json` where `XX` is the [BCP-47](https://tools.ietf.org/rfc/bcp/bcp47.txt) code of the new language. If in doubt how to choose the BCP-47 language code, use a two letter [ISO-631-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes). Only the `"translation":` line has to be translated, the `"defaultMessage"`, `"description"` etc should NOT be translated, they serve as help only, i.e.:

```js
"action_block_contact": {
  "translation": "Bloquear contacto", // <<<---- Only this string needs to be translated
  "defaultMessage": "Block Contact",  // This is the default message in English
  "description": "Flat button [Block Contact]", // This is an explanation where/how the string is used.
  "missing": false,
  "obsolete": false
},
```

When translating the `service-worker.js`, just add the strings directly to the file. Only two strings need to be translated "New message" and "New chat":

```js
const i18n = {
  ...
  'XX': {
    'new_message': "New message",
    'new_chat': "New chat",
  },
  ...
```

Please send a pull request with the new files. If you don't know how to create a pull request just sent the files in any way you can.


## Android

A single file needs to be translated: [/tinode/tindroid/app/src/main/res/values/strings.xml](https://github.com/tinode/tindroid/blob/devel/app/src/main/res/values/strings.xml)

Create a new directory `values-XX` in [app/src/main/res](https://github.com/tinode/tindroid/tree/devel/app/src/main/res), where `XX` is a two-letter [ISO-631-1](https://en.wikipedia.org/wiki/List_of_ISO_639-1_codes) code , optionally followed by a two letter [ISO 3166-1-alpha-2](https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2) region code (preceded by lowercase r). For example `values-pt.xml` would contain Portuguese translations while `values-pt-rBR.xml` is for _Brazilian_ Portuguese translations.

Make a copy of the file with English strings, place it to the new directory. Translate all the strings not marked with `translatable="false"` (the strings with `translatable="false"` don't need to be included at all) then send a pull request with the new file. If you don't know how to create a pull request then just sent the file in any way you can.


## iOS

Unfortunately iOS localization process is very convoluted and generally requires `Xcode` which runs only on Mac. Unless you are familiar with the iOS development, please create a [feature request](https://github.com/tinode/ios/issues/new?assignees=&labels=&template=feature_request.md&title=) for the desired language and we will send you the file for translation.

If you feel brave enough, you can translate the [following .xliff file](https://github.com/tinode/ios/blob/devel/Localizations/en.xcloc/Localized%20Contents/en.xliff) and then send it to us in any way you can. Translate all the strings between the `<target>` tags:

```xml
<trans-unit id="Action failed: %@" xml:space="preserve"> <!-- Do NOT change this line -->
  <source>Action failed: %@</source> <!-- This is the default message in English. -->
  <target>Se ha producido un error al realizar la acción: %@</target> <!-- Only this string "target" needs to be translated. -->
  <note>Toast notification</note> <!-- This is an explanation where/how the string is used. -->
</trans-unit>
```

If you are familiar with Xcode and localization for iOS, the exported localizations are located at [/Localizations](https://github.com/tinode/ios/tree/devel/Localizations).
