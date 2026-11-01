# Mayday HD

<img width="1920" height="1080" alt="Zrzut ekranu (2)" src="https://github.com/user-attachments/assets/f7d9a624-a963-4512-9873-caae58f2ec2e" />

## Co to jest?

[English summary below](#english)

**Mayday HD 1.0.1 by palin** – nieoficjalny patch fanowski do polskiej wersji gry
**Mayday: Conflict Earth (1.1p)**, który pozwala grać na Windows 10 i 11 w 1920×1080,
z płynnym obrazem i kilkoma poprawkami rozgrywki.

> Repozytorium **nie zawiera gry ani żadnych jej plików**. Do instalacji potrzebna jest
> oryginalna polska płyta CD albo jej pełny obraz (MDF/MDS lub BIN/CUE).

[English summary below](#english)

## Pobieranie

Gotowy patcher (`MaydayHD_1.0.1.zip`) znajdziesz w zakładce **[Releases](../../releases)**.

Windows może ostrzec, że plik pochodzi z internetu (SmartScreen). Kliknij
**Więcej informacji → Uruchom mimo to**. Sumę kontrolną pliku podaje `SHA256.txt` w archiwum.

## Instalacja

Proces krok po kroku w wiki: [Instalacja](https://github.com/palin/mayday-conflict-earth-hd/wiki/Instalacja)

## Zmiany w 1.0.1

- Naprawione klawisze w grze przy odświeżaniu powyżej 60 Hz: pauza (ESC), cele misji (M), opcje (O), kapitulacja, grupy bojowe (CTRL+cyfra) i tryby walki (F1–F4)

Masz już wersję 1.0? Uruchom `MaydayHD_1.0.1.exe` i kliknij **Zaktualizuj do 1.0.1**.

## Zawartość patcha 1.0.1

**Grafika**
- Rozdzielczość 1920×1080 - ustawiona jako domyślna (zastępuje 1024×768)
- Płynniejszy obraz - zgodny z odświeżaniem monitora- np. 144Hz (prędkość gry bez zmian)
- Płynniejsze przewijanie mapy

**Minimapa w grze**
- Lepsza przejrzystość i orientacja (jednolite tło, obrys świata, jaśniejszy teren, obwódki jednostek)
- Dokładna ramka widoku (kliknięcie na minimapie precyzyjnie centruje widok)

**Rozgrywka**
- Poprawione zaznaczanie jednostek (wciąż mogą być jednak problemy przy zaznaczaniu jednostek na różnych wysokościach terenu)

**Terminal gry**
- Spokojniejsze przewijanie widoku prawo-lewo

**Windows 10/11**
- Dodanie narzędzia "cnc-ddraw" (wymagane, by gra działała) + program „cnc-ddraw config” (w folderze z grą), gdzie można dostosować sobie mnóstwo opcji związanych z wyświetlaniem gry
- Naprawa potencjalnego problemu z wyciszonym dźwiękiem gry
- Instalacja kodeka Power!Video z płyty gry – działają filmy w intro i odprawach misji (Windows jednorazowo zapyta o zgodę administratora)

**Intro**
- Domyślnie wyłączone aby przyspieszyć uruchamianie gry (opcja w patcherze do włączenia intro)

## Instalacja w skrócie

1. Zainstaluj [otvdm (winevdm)](https://github.com/otya128/winevdm/releases) – pozwala uruchomić 16-bitowy instalator gry.
2. Zamontuj płytę lub obraz (np. DAEMON Tools Lite; plik `.iso` nie wystarczy – brak w nim ścieżek z muzyką).
3. W PowerShellu uruchomionym jako administrator uruchom `SETUP.EXE` z płyty (nie `MDStart.exe`).
   Wybierz instalację maksymalną i folder poza „Program Files”. Instalację Power!Video i DirectX anuluj.
4. Uruchom `MaydayHD_1.0.1.exe` (najlepiej z folderu gry), z zamontowaną płytą, i kliknij **Zainstaluj patch**.

Szczegółowa instrukcja krok po kroku: [docs/README_patch_PL.txt](docs/README_patch_PL.txt).

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
go build -buildvcs=false -trimpath -ldflags "-H=windowsgui -s -w" -o MaydayHD_1.0.1.exe .
```

`go test ./...` uruchamia testy silnika patcha; testy wymagające plików gry są pomijane, jeśli ich nie ma.

## Licencja

Kod patchera i źródła poprawek: [MIT](LICENSE).
cnc-ddraw: licencja MIT, © FunkyFr3sh.
Mayday: Conflict Earth jest własnością jej wydawców i twórców. To nieoficjalny projekt fanowski,
niezwiązany z właścicielami praw do gry.

---

## English

**Mayday HD 1.0.1 by palin** is an unofficial fan patch for the Polish release of
*Mayday: Conflict Earth* (version 1.1p) on Windows 10/11. 
The patch will work only on the Polish version of the game. If you'd like to try it on a different version, like English, German or other, I'd need to create a separate patch for it. Please raise an issue at the top and let me know!

Changes in 1.0.1:
- 1920×1080 by default, with the frame rate matching the monitor (e.g. 144 Hz) and unchanged game speed
- smooth map scrolling
- minimap improvements with an accurate view frame
- better unit selection and calmer terminal scrolling
- cnc-ddraw (required on Windows 10/11) and a sound fix
- the Power!Video codec for the intro and briefing videos, extracted from the player's own CD
- intro off by default

The repository contains **no game files**. You need the original Polish CD or a full image of it.
Download the patcher from **[Releases](../../releases)**, run it from the game folder with the CD mounted
and click **Zainstaluj patch** (Install). The original files are kept in `MaydayHD_backup`
and can be restored with **Przywróć oryginał** (Restore original).
