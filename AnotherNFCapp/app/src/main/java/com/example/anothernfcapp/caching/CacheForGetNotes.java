package com.example.anothernfcapp.caching;

import android.content.Context;
import android.util.Log;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.FileReader;
import java.io.FileWriter;
import java.io.IOException;

public class CacheForGetNotes {
    private Context context;
    private File file;
    private BufferedWriter bufferedWriter;
    private BufferedReader bufferedReader;

    public CacheForGetNotes(Context context){
        this.context = context;
        file = new File(context.getCacheDir(), "cacheForReceivedNotes");
    }

    
    public void writeToCache(String data) throws IOException {
        bufferedWriter = new BufferedWriter(new FileWriter(file));
        Log.d("CACHEGET", data);
        try {
            bufferedWriter.append(data).append("\n");
            Log.d("CACHEGET", "Added your notes to the cache");
            bufferedWriter.close();
        } catch (IOException e) {
            e.printStackTrace();
            Log.e("CACHEGET", "IOEXCEPTION");
        }
    }
    
    public void clearCache(){
        try {
            bufferedWriter = new BufferedWriter(new FileWriter(file));
        } catch (IOException e) {
            Log.e("CACHEGET", e.toString());
        }
        try {
            bufferedWriter.write("");
            Log.d("CACHEGET", "Cache was cleared");
            bufferedWriter.close();
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
    
    public String getCachedNotes() throws IOException {
        StringBuilder stringBuilder = new StringBuilder();
        bufferedReader = new BufferedReader(new FileReader(file));
        String line;
        while ((line = bufferedReader.readLine()) != null){
            stringBuilder.append(line);
        }
        return stringBuilder.toString();
    }


}
