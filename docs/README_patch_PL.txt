MAYDAY HD 1.0.2 by palin
dla Mayday: Conflict Earth PL (1.1p)
=======================================

[PL]

Zmiany w 1.0.2
• Lista zapisów (ODCZYT / ZAPIS): strzałki przewijają o jedną pozycję na kliknięcie, przytrzymanie przewija płynnie
• Lista zapisów otwiera się na końcu – od razu widać najnowsze zapisy
• Lista zapisów posortowana według daty – najnowszy zapis jest zawsze na dole
• Szybkie kliknięcia nie są już gubione (np. strzałki w monitorze badań, wydawanie rozkazów)
• Panel WIDEO pokazuje jednostkę latającą także po zaznaczeniu ramką (wcześniej pusty podgląd, a czasem nagłe zamknięcie gry)

Zmiany w 1.0.1
• Naprawione klawisze w grze przy odświeżaniu powyżej 60 Hz: pauza (ESC), cele misji (M), opcje (O), kapitulacja, grupy bojowe (CTRL+cyfra) i tryby walki (F1–F4)

Zawartość patcha 1.0.2:

Grafika
• Rozdzielczość 1920×1080 - ustawiona jako domyślna (zastępuje 1024×768)
• Płynniejszy obraz - zgodny z odświeżaniem monitora- np. 144Hz (prędkość gry bez zmian)
• Płynniejsze przewijanie mapy

Minimapa w grze
• Lepsza przejrzystość i orientacja (jednolite tło, obrys świata, jaśniejszy teren, obwódki jednostek)
• Dokładna ramka widoku (kliknięcie na minimapie precyzyjnie centruje widok)

Rozgrywka
• Poprawione zaznaczanie jednostek (wciąż mogą być jednak problemy przy zaznaczaniu jednostek na różnych wysokościach terenu)

Terminal gry
• Spokojniejsze przewijanie widoku prawo-lewo

Windows 10/11
• Dodanie narzędzia "cnc-ddraw" (wymagane, by gra działała) + program „cnc-ddraw config” (w folderze z grą), gdzie można dostosować sobie mnóstwo opcji związanych z wyświetlaniem gry
• Naprawa potencjalnego problemu z wyciszonym dźwiękiem gry
• Instalacja kodeka Power!Video z płyty gry – działają filmy w intro i odprawach misji (Windows jednorazowo zapyta o zgodę administratora)

Intro
• Domyślnie wyłączone aby przyspieszyć uruchamianie gry (opcja poniżej do włączenia intro)

Jak zainstalować
----------------
1. Zainstaluj grę z CD: otwórz PowerShell jako administrator, przejdź na napęd
   z płytą (np. F:) i uruchom .\SETUP.EXE (nie MDStart.exe). Wybierz folder
   poza „Program Files”, np. C:\Gry\Mayday. Instalację Power!Video i DirectX 5
   możesz anulować.
2. Zamknij grę. Uruchom MaydayHD_1.0.2.exe - najlepiej z folderu gry, wtedy sam
   go znajdzie; inaczej wskaż folder przyciskiem „Wybierz…”.
3. Jeśli chcesz intro, zaznacz „Odtwarzaj intro przy starcie”.
4. Kliknij „Zainstaluj patch”.

Masz już wersję 1.0 lub 1.0.1? Uruchom MaydayHD_1.0.2.exe i kliknij „Zaktualizuj do 1.0.2”.

Patcher można uruchomić ponownie w każdej chwili: zmienić ustawienie intro albo
przywrócić oryginalne pliki („Przywróć oryginał”). Nowsza wersja patcha
rozpozna starszą i zaproponuje aktualizację.

Patcher zmienia Mayday.exe, mayday.art i rozdzielczość w mayday.cfg, instaluje
cnc-ddraw (ddraw.dll, ddraw.ini, folder Shaders, „cnc-ddraw config.exe”),
kodek filmów Power!Video wypakowany z płyty gracza (PV_WIN\SETUP_PV.EXE →
pvez32.drv i Pvezw32d.dll w C:\Windows\SysWOW64 + wpis vidc.pvez w rejestrze)
i zmienia nazwę Intro.exe, gdy intro jest wyłączone. Oryginały zapisuje
w folderze MaydayHD_backup. Działa tylko z polską wersją 1.1p - inne wersje
odrzuca bez zmian.

Wskazówki
---------
- Windows może ostrzec przed programem z internetu (SmartScreen):
  „Więcej informacji” → „Uruchom mimo to”. Sumę kontrolną pliku podaje SHA256.txt.
- Gra nadal wymaga zamontowanego obrazu płyty (np. DAEMON Tools, plik .cue).
- Brak dźwięku na starcie: uruchom grę, ustaw Mayday na 100 w Mikserze
  głośności Windows i wyjdź z gry przez menu. Wystarczy raz.
- Podczas instalacji patcha płyta (obraz) musi być zamontowana - z niej
  instalowany jest kodek filmów; Windows zapyta wtedy raz o zgodę administratora.
- cnc-ddraw pochodzi z github.com/FunkyFr3sh/cnc-ddraw (licencja MIT,
  plik cnc-ddraw-LICENSE.txt).

[EN]

Mayday HD 1.0.2 by palin - unofficial fan patch for the Polish Mayday: Conflict
Earth 1.1p on Windows 10/11: 1920x1080 by default, frame rate matching the
monitor (e.g. 144 Hz) with unchanged game speed, smooth map scrolling, minimap
improvements with an accurate view frame, better unit selection, calmer
terminal scrolling, cnc-ddraw (required on Windows 10/11), sound fix, the Power!Video codec for
the intro and briefing videos (extracted from the player's own CD), intro off
by default. Run MaydayHD_1.0.2.exe, choose the game folder and click
"Zainstaluj patch" (Install). "Przywróć oryginał" (Restore original) puts the
original files back; they are kept in MaydayHD_backup. Only the Polish 1.1p
files are accepted. The game CD image is still required.

Checksums (SHA256) of the supported original files:
  Mayday.exe  9F9AB3F11A25BB2241537E984E1E487EB8E97282DBEB20F6AB0D137F359EBFB4
  mayday.art  08016DE9343D6FD0971FF6D9ECE4904D33BF7D26E1FE6502A7EDD6B0D347C613
