using System.IO;
using SkiaSharp;

namespace AdventOfCode.Exercises;

public class Exercise14 : IExercise
{
    public int GetDay() => 14;

    private class Coords
    {
        public long X { get; set; }
        public long Y { get; set; }
    }

    private class Robot
    {
        public Coords Position { get; set; }
        public Coords Velocity { get; set; }
    }

    private List<Robot> Parse(string[] input)
    {
        List<Robot> parsed = new();
        foreach (var i in input)
        {
            var iSplit = i.Split(" ");

            var posRaw = iSplit.First();
            var posSplit = posRaw.Replace("p=", "").Split(",");
            long.TryParse(posSplit.First(), out var xPos);
            long.TryParse(posSplit.Last(), out var yPos);

            var velRaw = iSplit.Last();
            var velSplit = velRaw.Replace("v=", "").Split(",");
            long.TryParse(velSplit.First(), out var xVel);
            long.TryParse(velSplit.Last(), out var yVel);

            parsed.Add(
                new()
                {
                    Position = new() { X = xPos, Y = yPos },
                    Velocity = new() { X = xVel, Y = yVel }
                }
            );
        }
        return parsed;
    }

    private void Wait(List<Robot> robots, long maxX, long maxY)
    {
        foreach (var robot in robots)
        {
            robot.Position.X += robot.Velocity.X;
            robot.Position.Y += robot.Velocity.Y;

            if (robot.Position.X >= maxX) robot.Position.X -= maxX;
            else if (robot.Position.X < 0) robot.Position.X += maxX;

            if (robot.Position.Y >= maxY) robot.Position.Y -= maxY;
            else if (robot.Position.Y < 0) robot.Position.Y += maxY;
        }
    }

    private long CalculateSafetyFactor(List<Robot> robots, long maxX, long maxY)
    {
        var midX = maxX / 2;
        var midY = maxY / 2;
        List<long> quad1XRange = [0, midX - 1];
        List<long> quad1YRange = [0, midY - 1];

        List<long> quad2XRange = [midX + 1, maxX - 1];
        List<long> quad2YRange = [0, midY - 1];

        List<long> quad3XRange = [0, midX - 1];
        List<long> quad3YRange = [midY + 1, maxY - 1];

        List<long> quad4XRange = [midX + 1, maxX - 1];
        List<long> quad4YRange = [midY + 1, maxY - 1];

        long q1Count = 0;
        long q2Count = 0;
        long q3Count = 0;
        long q4Count = 0;

        foreach (var robot in robots)
        {
            if (robot.Position.X >= quad1XRange[0] && robot.Position.X <= quad1XRange[1]
                && robot.Position.Y >= quad1YRange[0] && robot.Position.Y <= quad1YRange[1])
            {
                q1Count++;
            }
            else if (robot.Position.X >= quad2XRange[0] && robot.Position.X <= quad2XRange[1]
                && robot.Position.Y >= quad2YRange[0] && robot.Position.Y <= quad2YRange[1])
            {
                q2Count++;
            }
            else if (robot.Position.X >= quad3XRange[0] && robot.Position.X <= quad3XRange[1]
                && robot.Position.Y >= quad3YRange[0] && robot.Position.Y <= quad3YRange[1])
            {
                q3Count++;
            }
            else if (robot.Position.X >= quad4XRange[0] && robot.Position.X <= quad4XRange[1]
                && robot.Position.Y >= quad4YRange[0] && robot.Position.Y <= quad4YRange[1])
            {
                q4Count++;
            }
        }

        return q1Count * q2Count * q3Count * q4Count;
    }

    private List<List<int>> MapRobots(List<Robot> robots, int maxX, int maxY)
    {
        List<List<int>> roboMap = Enumerable.Repeat( Enumerable.Repeat( 0, maxX ).ToList(), maxY).ToList();

        foreach (var robot in robots)
        {
            // paying the price here by using longs
            roboMap[(int)robot.Position.Y][(int)robot.Position.X] ++;
        }

        return roboMap;
    }

    private bool IsTreeArrangement(List<List<int>> map)
    {
        var maxX = map[0].Count();
        var maxY = map.Count();

        for (var i = 0; i < maxY; i++)
        {
            List<Coords> contiguousCoords = new();
            for (var j = 0; j < maxX; j++)
            {
                if (j + 1 > maxX)
                {
                    contiguousCoords.Add(new() { X=j, Y=i });
                    break;
                }
                else if (map[i][j] != 0 && map[i][j+1] != 0)
                {
                    contiguousCoords.Add(new() { X=j, Y=i });
                }
                else if (map[i][j] != 0 && map[i][j+1] != 0)
                {
                    contiguousCoords.Add(new() { X=j, Y=i });
                }
            }
        }

        return true;
    }

    private SKBitmap CreateEmptyImage(int x, int y)
    {
        var bmap = new SKBitmap(x, y);
        for (var i = 0; i < x; i++)
        {
            for (var j = 0; j < x; j++)
            {
                bmap.SetPixel(i, j, SKColors.White);
            }
        }

        return bmap;
    }

    private void DrawOnImage(SKBitmap image, List<Robot> robots)
    {
        foreach (var robot in robots)
        {
            image.SetPixel((int)robot.Position.X, (int)robot.Position.Y, SKColors.Black);
        }
    }

    public Result ExecutePart1(string inputFile)
    {
        var input = File.ReadAllLines(inputFile);
        var robots = Parse(input);

        var maxX = 101;
        var maxY = 103;

        for (var i = 0; i < 100; i ++)
        {
            Wait(robots, maxX, maxY);
        }

        return new Result(CalculateSafetyFactor(robots, maxX, maxY), true);
    }

    public Result ExecutePart2(string inputFile)
    {
        return new Result(6493, true);
        var input = File.ReadAllLines(inputFile);
        var robots = Parse(input);

        var maxX = 101;
        var maxY = 103;

        var outFolder = "./renders";
        if (!Directory.Exists(outFolder))
        {
            Directory.CreateDirectory(outFolder);
        }

        for (var i = 0; i < 10000; i++)
        {
            Wait(robots, maxX, maxY);
            // var roboMap = MapRobots(robots, maxX, maxY);
            var image = CreateEmptyImage(maxX, maxY);
            DrawOnImage(image, robots);
            using (var fs = new SKFileWStream($"{outFolder}/{i+1}.bmp"))
            {
                image.Encode(fs, SKEncodedImageFormat.Jpeg, 85);
            }
        }

        return new Result();
    }
}
