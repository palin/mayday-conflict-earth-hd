# Mayday: Conflict Earth HD Patch

<img width="1920" height="1080" alt="Mayday: Conflict Earth gameplay screenshot" src="https://github.com/user-attachments/assets/f7d9a624-a963-4512-9873-caae58f2ec2e" />

## Co to jest?

[English summary at the bottom](#english)

**Mayday HD by palin** – nieoficjalny patch fanowski do polskiej wersji gry
**Mayday: Conflict Earth (1.1p)** (Mayday: Dzień Rozstrzygnięcia), który pozwala grać na Windows 10 i 11 w 1920×1080 (Full HD),
z płynnym obrazem i poprawkami rozgrywki.

> Repozytorium **nie zawiera gry ani żadnych jej plików**. Do instalacji potrzebna jest
> oryginalna polska płyta CD albo jej pełny obraz (MDF/MDS lub BIN/CUE).

## Pobieranie

Gotowy patcher (`MaydayHD_x.y.z.zip`) znajdziesz w zakładce **[Releases](../../releases)**.

Windows może ostrzec, że plik pochodzi z internetu (SmartScreen). Kliknij
**Więcej informacji → Uruchom mimo to**. Sumę kontrolną pliku podaje `SHA256.txt` w archiwum.

## Instalacja

Proces krok po kroku w wiki: [Instalacja](https://github.com/palin/mayday-conflict-earth-hd/wiki/Instalacja)

## Zmiany

Wszystkie zmiany są podane w [changelog.md](https://github.com/palin/mayday-conflict-earth-hd/blob/main/CHANGELOG.md)

## Cała obecna zawartość patcha

**Grafika**
- Rozdzielczość 1920×1080 - ustawiona jako domyślna (zastępuje 1024×768)
- Płynniejszy obraz - zgodny z odświeżaniem monitora- np. 144Hz (prędkość gry bez zmian)
- Płynniejsze przewijanie mapy

**Minimapa w grze**
- Lepsza przejrzystość i orientacja (jednolite tło, obrys świata, jaśniejszy teren, obwódki jednostek)
- Dokładna ramka widoku (kliknięcie na minimapie precyzyjnie centruje widok)

**Rozgrywka**
- Poprawione zaznaczanie jednostek (wciąż mogą być jednak problemy przy zaznaczaniu jednostek na różnych wysokościach terenu)
- Działające klawisze przy odświeżaniu powyżej 60 Hz: pauza (ESC), cele misji (M), opcje (O), kapitulacja, grupy bojowe (CTRL+cyfra) i tryby walki (F1–F4)
- Panel WIDEO pokazuje jednostkę latającą także po zaznaczeniu ramką (wcześniej pusty podgląd, a czasem nagłe zamknięcie gry)
- Szybkie kliknięcia nie są już gubione (np. strzałki w monitorze badań, wydawanie rozkazów)

**Menu**
- Lista zapisów (ODCZYT / ZAPIS): strzałki przewijają o jedną pozycję na kliknięcie, przytrzymanie przewija płynnie
- Lista zapisów posortowana według daty i otwiera się na końcu – najnowszy zapis jest zawsze na dole

**Terminal gry**
- Spokojniejsze przewijanie widoku prawo-lewo

**Windows 10/11**
- Dodanie narzędzia "cnc-ddraw" (wymagane, by gra działała) + program „cnc-ddraw config” (w folderze z grą), gdzie można dostosować sobie mnóstwo opcji związanych z wyświetlaniem gry
- Naprawa potencjalnego problemu z wyciszonym dźwiękiem gry
- Instalacja kodeka Power!Video z płyty gry – działają filmy w intro i odprawach misji (Windows jednorazowo zapyta o zgodę administratora)

**Intro**
- Domyślnie wyłączone aby przyspieszyć uruchamianie gry (opcja w patcherze do włączenia intro)

## Instalacja w skrócie

Szczegółowa instrukcja krok po kroku w Wiki: [Instalacja](https://github.com/palin/mayday-conflict-earth-hd/wiki/Instalacja).

1. Zainstaluj [otvdm (winevdm)](https://github.com/otya128/winevdm/releases) – pozwala uruchomić 16-bitowy instalator gry.
2. Zamontuj płytę lub obraz (np. DAEMON Tools Lite; plik `.iso` nie wystarczy – brak w nim ścieżek z muzyką).
3. W PowerShellu uruchomionym jako administrator uruchom `SETUP.EXE` z płyty (nie uruchamiaj `MDStart.exe` ani `SETUP.exe` bezpośrednio, musi być przez PowerShell).
   Wybierz instalację maksymalną i folder poza „Program Files”. Instalację Power!Video i DirectX anuluj.
4. Pobierz i uruchom najnowszy patch `MaydayHD_x.y.z.exe` (najlepiej z folderu gry), z zamontowaną płytą, i kliknij **Zainstaluj patch**.

Patcher sprawdza sumy kontrolne plików (działa tylko z polską wersją 1.1p), zapisuje oryginały
w folderze `MaydayHD_backup` i potrafi je przywrócić przyciskiem **Przywróć oryginał**.

| Plik | SHA256 oryginału |
| --- | --- |
| Mayday.exe | `9F9AB3F11A25BB2241537E984E1E487EB8E97282DBEB20F6AB0D137F359EBFB4` |
| mayday.art | `08016DE9343D6FD0971FF6D9ECE4904D33BF7D26E1FE6502A7EDD6B0D347C613` |

## Zawartość repozytorium

| Folder | Co zawiera |
| --- | --- |
| `patcher/` | Kod patchera (Go, tylko biblioteka standardowa, interfejs Win32) |
| `patcher/ddraw/` | Dołączony [cnc-ddraw](https://github.com/FunkyFr3sh/cnc-ddraw) (licencja MIT, `cnc-ddraw-LICENSE.txt`) |
| `asm/` | Źródła asemblerowe poprawek w `Mayday.exe` (60/144 kl./s, płynne przewijanie, ramka minimapy, dźwięk, terminal) |
| `asm/scripts/` | Skrypty Python, którymi powstawały kolejne poprawki (dokumentacja zmian) |
| `docs/` | README dołączane do patcha |

Kodek Power!Video **nie jest** częścią repozytorium – patcher wypakowuje go z pliku
`PV_WIN\SETUP_PV.EXE` na płycie gracza.

## Budowanie

Potrzebny jest [Go](https://go.dev/dl/) 1.24 lub nowszy. W PowerShellu:

```powershell
cd patcher
$env:GOOS="windows"; $env:GOARCH="amd64"
go build -buildvcs=false -trimpath -ldflags "-H=windowsgui -s -w" -o MaydayHD_1.0.2.exe .
```

`go test ./...` uruchamia testy silnika patcha; testy wymagające plików gry są pomijane, jeśli ich nie ma.

## Buy me a coffee!

Jeśli masz ochotę docenić moją pracę nad patchem, postaw mi kawę:

[<img width="132" height="67" alt="logotyp-pionowy-roz" src="https://github.com/user-attachments/assets/f200af0b-629a-4957-a5a4-a4cb346e9e51" />](https://buycoffee.to/palin)

[Postaw mi kawę / Buy me a coffee](https://buycoffee.to/palin)

## Licencja

Kod patchera i źródła poprawek: [MIT](LICENSE).
cnc-ddraw: licencja MIT, © FunkyFr3sh.
Mayday: Conflict Earth jest własnością jej wydawców i twórców. To nieoficjalny projekt fanowski,
niezwiązany z właścicielami praw do gry.

---

## English

**Mayday HD by palin** is an unofficial fan patch for the Polish release of
*Mayday: Conflict Earth* (version 1.1p) on Windows 10/11. 
The patch will work only on the Polish version of the game. If you'd like to try it on a different version, like English, German or other, I'd need to create a separate patch for it. Please raise an issue at the top and let me know!

Included in the patch:

- 1920×1080 by default, with the frame rate matching the monitor (e.g. 144 Hz) and unchanged game speed
- smooth map scrolling
- minimap improvements with an accurate view frame
- better unit selection and calmer terminal scrolling
- cnc-ddraw (required on Windows 10/11) and a sound fix
- the Power!Video codec for the intro and briefing videos, extracted from the player's own CD
- intro off by default
- the save list arrows scroll one line per click 
- the save list is sorted by date and opens at the newest saves 
- fast clicks are no longer lost 
- the VIDEO panel shows a flying unit after a box selection (before: blank thumbnail, sometimes a crash)

The repository contains **no game files**. You need the original Polish CD or a full image of it.
Download the patcher from **[Releases](../../releases)**, run it from the game folder with the CD mounted
and click **Zainstaluj patch** (Install). The original files are kept in `MaydayHD_backup`
and can be restored with **Przywróć oryginał** (Restore original).
