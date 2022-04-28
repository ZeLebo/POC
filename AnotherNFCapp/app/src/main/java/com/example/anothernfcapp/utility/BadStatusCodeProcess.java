package com.example.anothernfcapp.utility;

import android.content.Context;
import android.widget.Toast;

public class BadStatusCodeProcess {
    public static void parseBadStatusCode(int statusCode, String response, Context context){
        if (statusCode == 0){
            Toast.makeText(context, "No internet connection", Toast.LENGTH_SHORT).show();
        }
        else if (statusCode >= 400){
            if (response.equals("")){
                Toast.makeText(context, "Error " + statusCode, Toast.LENGTH_SHORT).show();
            }
            else{
                Toast.makeText(context, "Error " + statusCode + ": " + response, Toast.LENGTH_SHORT).show();
            }
        }
    }
}