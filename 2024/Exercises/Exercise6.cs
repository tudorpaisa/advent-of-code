using AdventOfCode.Utilities;

namespace AdventOfCode.Exercises;

public class Exercise6 : IExercise
{
    public int GetDay() => 6;

    private int[][] _directions = new int[][]
    {
        new int[] {-1, 0 },
        new int[] { 0, 1 },
        new int[] { 1, 0 },
        new int[] { 0, -1},
    };
    private char[] _directionChars = new char[] { '^', '>', 'V', '<' };

    private int GetNextDirectionIndex(int idx) => (idx + 1) % 4;

    private int[] FindStartPos(string[] input)
    {
        for (var i = 0; i < input.Count(); i++)
        {
            var row = input[i].ToList();

            foreach (var dir in _directionChars)
            {
                if (row.Contains(dir)) return new int[] { i, row.IndexOf(dir) };
            }
        }
        return new int[2];
    }

    private int[] GetNewPosition(int[] pos, int dirIdx)
    {
        var dirDelta = _directions[dirIdx];
        var newPos = new int[2];
        newPos[0] = dirDelta[0] + pos[0];
        newPos[1] = dirDelta[1] + pos[1];
        return newPos;
    }

    private int CountPositions(string[] input, int[] pos, int directionIdx)
    {
        var visited = input.Select(r => Enumerable.Repeat(false, r.Count()).ToArray()).ToArray();
        visited[pos[0]][pos[1]] = true;
        var count = 1;
        while (true)
        {
            var newPos = GetNewPosition(pos, directionIdx);

            if (newPos[0] >= input.Count() || newPos[1] >= input[0].Count() || newPos[0] < 0 || newPos[1] < 0)
            {
                break;
            }
            else if (input[newPos[0]][newPos[1]] == '#')
            {
                directionIdx = GetNextDirectionIndex(directionIdx);
            }
            else
            {
                pos = newPos;
                count++;
                visited[pos[0]][pos[1]] = true;
            }
        }

        return visited.SelectMany(r => r).Where(b => b == true).Count();
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);

        var startPos = FindStartPos(input);
        var directionIdx = 0;


        return new Result(CountPositions(input, startPos, directionIdx), true);
    }

    private int[][] GetTileVisitationCount(string[] input, int[] pos, int directionIdx)
    {
        var visited = input.Select(r => Enumerable.Repeat(0, r.Count()).ToArray()).ToArray();
        visited[pos[0]][pos[1]] = 1;
        while (true)
        {
            var newPos = GetNewPosition(pos, directionIdx);

            if (newPos[0] >= input.Count() || newPos[1] >= input[0].Count() || newPos[0] < 0 || newPos[1] < 0)
            {
                break;
            }
            else if (input[newPos[0]][newPos[1]] == '#')
            {
                directionIdx = GetNextDirectionIndex(directionIdx);
            }
            else
            {
                pos = newPos;
                visited[pos[0]][pos[1]]++;
            }
        }

        return visited;
    }

    private int[][] GetVisitedTileLocations(int[][] count)
    {
        List<int[]> results = [];
        for (var i = 0; i < count.Count(); i++)
        {
            var multipleVisitTiles = count[i].Select(( n, idx ) => (n, idx)).Where(n => n.n >= 1).Select(n => n.idx);
            foreach( var j in multipleVisitTiles )
            {
                results.Add(new int[] { i, j });
            }
        }
        return results.ToArray();
    }

    private int[][] GetPlacesToBlock(string[] input, int[] pos)
    {
        List<int[]> results = new();

        foreach (var d in _directions)
        {
            var newR = pos[0] + d[0];
            var newC = pos[1] + d[1];

            if (newR >= input.Count() || newC >= input.Count() || newR < 0 || newC < 0) continue;
            if (input[newR][newC] == '#') continue;

            results.Add(new int[] { newR, newC });
        }

        return results.ToArray();
    }

    private bool CheckIfLoop(string[] input, int[] pos, int directionIdx)
    {
        Dictionary<string, bool> visitMap = new();
        while (true)
        {
            var newPos = GetNewPosition(pos, directionIdx);

            if (newPos[0] >= input.Count() || newPos[1] >= input[0].Count() || newPos[0] < 0 || newPos[1] < 0)
            {
                break;
            }
            else if (input[newPos[0]][newPos[1]] == '#')
            {
                directionIdx = GetNextDirectionIndex(directionIdx);
            }
            else
            {
                pos = newPos;

                var key = $"{pos[0]} {pos[1]} {directionIdx}";
                if (visitMap.ContainsKey(key))
                {
                    return true;
                }
                visitMap[key] = true;
            }
        }

        return false;
    }

    public Result ExecutePart2(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var startPos = FindStartPos(input);
        var directionIdx = 0;
        var visCount = GetTileVisitationCount(input, startPos, directionIdx);
        var crossedCoords = GetVisitedTileLocations(visCount);

        var count = 0;
        Dictionary<string, bool> blockTries = new();

        foreach (var blockCoord in crossedCoords)
        {
            if (blockCoord[0] == startPos[0] && blockCoord[1] == startPos[1]) continue;
            var key = $"{blockCoord[0]} {blockCoord[1]}";
            if (blockTries.ContainsKey(key)) continue;
            blockTries[key] = true;
            var oldRow = input[blockCoord[0]];
            var newRow = oldRow[0..blockCoord[1]] + '#' + oldRow[(blockCoord[1] + 1)..oldRow.Count()];
            input[blockCoord[0]] = newRow;
            if (CheckIfLoop(input, startPos, directionIdx)) count++;
            input[blockCoord[0]] = oldRow;
        }

        return new Result(count, true);
    }
}
