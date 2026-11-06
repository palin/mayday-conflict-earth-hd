# Changelog – Mayday HD by palin

Nieoficjalny patch fanowski do polskiej wersji *Mayday: Conflict Earth* 1.1p (Windows 10/11).
Pobieranie: [Releases](https://github.com/palin/mayday-conflict-earth-hd/releases).
Każda wersja patchera rozpoznaje starsze i proponuje aktualizację („Zaktualizuj do …”).

*Unofficial fan patch for the Polish Mayday: Conflict Earth 1.1p. English summary under each version.*

---

## [1.0.2] – 2026-11-06

### Menu
- Lista zapisów (ODCZYT / ZAPIS): strzałki przewijają o jedną pozycję na kliknięcie, przytrzymanie przewija płynnie
- Lista zapisów otwiera się na końcu – od razu widać najnowsze zapisy
- Lista zapisów posortowana według daty – najnowszy zapis jest zawsze na dole
- Szybkie kliknięcia nie są już gubione (np. strzałki w monitorze badań, wydawanie rozkazów)

### Rozgrywka
- Panel WIDEO pokazuje jednostkę latającą także po zaznaczeniu ramką (wcześniej pusty podgląd, a dosyć często nagłe zamknięcie gry)

**EN:** Save list arrows scroll one line per click (smooth when held); the save list is sorted by
date and opens at the newest saves; fast clicks are no longer lost; the VIDEO panel shows a flying
unit after a box selection (before: blank thumbnail, sometimes a crash).

## [1.0.1] – 2026-11-01

### Poprawki
- Naprawione klawisze w grze przy odświeżaniu powyżej 60 Hz: pauza (ESC), cele misji (M), opcje (O), kapitulacja, grupy bojowe (CTRL+cyfra) i tryby walki (F1–F4)

**EN:** Fixed in-game keys on monitors above 60 Hz: pause (ESC), mission objectives (M),
options (O), surrender, battle groups (CTRL+number) and combat modes (F1–F4).

## [1.0] – 2026-11-01

Pierwsze wydanie.

### Grafika
- Rozdzielczość 1920×1080 - ustawiona jako domyślna (zastępuje 1024×768)
- Płynniejszy obraz - zgodny z odświeżaniem monitora- np. 144Hz (prędkość gry bez zmian)
- Płynniejsze przewijanie mapy

### Minimapa w grze
- Lepsza przejrzystość i orientacja (jednolite tło, obrys świata, jaśniejszy teren, obwódki jednostek)
- Dokładna ramka widoku (kliknięcie na minimapie precyzyjnie centruje widok)

### Rozgrywka
- Poprawione zaznaczanie jednostek (wciąż mogą być jednak problemy przy zaznaczaniu jednostek na różnych wysokościach terenu)

### Terminal gry
- Spokojniejsze przewijanie widoku prawo-lewo

### Windows 10/11
- Dodanie narzędzia "cnc-ddraw" (wymagane, by gra działała) + program „cnc-ddraw config” (w folderze z grą), gdzie można dostosować sobie mnóstwo opcji związanych z wyświetlaniem gry
- Naprawa potencjalnego problemu z wyciszonym dźwiękiem gry
- Instalacja kodeka Power!Video z płyty gry – działają filmy w intro i odprawach misji (Windows jednorazowo zapyta o zgodę administratora)

### Intro
- Domyślnie wyłączone aby przyspieszyć uruchamianie gry (opcja w patcherze do włączenia intro)

**EN:** First release: 1920×1080 by default, frame rate matching the monitor (e.g. 144 Hz) with
unchanged game speed, smooth map scrolling, clearer minimap with an accurate view frame, better
unit selection, calmer terminal scrolling, cnc-ddraw (required on Windows 10/11), muted-sound fix,
the Power!Video codec for intro and briefing videos (extracted from the player's own CD), intro off
by default.

[1.0.2]: https://github.com/palin/mayday-conflict-earth-hd/releases/tag/v1.0.2
[1.0.1]: https://github.com/palin/mayday-conflict-earth-hd/releases/tag/v1.0.1
[1.0]: https://github.com/palin/mayday-conflict-earth-hd/releases/tag/v1.0